package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthClient handles OAuth 2.0 PKCE authentication.
type OAuthClient struct {
	BaseURL      string
	HTTPClient   *http.Client
	Keyring      *Keyring
}

// TokenResponse represents the OAuth token endpoint response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// ClientRegistration represents the client registration response.
type ClientRegistration struct {
	ClientID string `json:"client_id"`
}

// NewOAuthClient creates a new OAuth client.
func NewOAuthClient(baseURL string, keyring *Keyring) *OAuthClient {
	return &OAuthClient{
		BaseURL:    strings.TrimSuffix(baseURL, "/"),
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Keyring:    keyring,
	}
}

// Login performs the OAuth PKCE flow.
func (c *OAuthClient) Login() (*TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: Register dynamic client
	clientID, err := c.registerClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("registering client: %w", err)
	}

	// Step 2: Generate PKCE parameters
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := generateState()

	// Step 3: Start local callback server
	callbackChan := make(chan string, 1)
	errorChan := make(chan error, 1)
	port, err := c.startCallbackServer(state, callbackChan, errorChan)
	if err != nil {
		return nil, fmt.Errorf("starting callback server: %w", err)
	}

	// Step 4: Build authorization URL and open browser
	authURL := c.buildAuthURL(clientID, codeChallenge, state, port)
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If the browser doesn't open, visit:\n  %s\n", authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Printf("(Could not open browser automatically: %v)\n", err)
	}

	// Step 5: Wait for callback
	select {
	case code := <-callbackChan:
		// Exchange code for tokens
		tokens, err := c.exchangeCode(ctx, clientID, code, codeVerifier, port)
		if err != nil {
			return nil, fmt.Errorf("exchanging code: %w", err)
		}

		// Store refresh token
		if err := c.Keyring.SaveRefreshToken(tokens.RefreshToken); err != nil {
			return nil, fmt.Errorf("storing refresh token: %w", err)
		}

		return tokens, nil

	case err := <-errorChan:
		return nil, err

	case <-ctx.Done():
		return nil, fmt.Errorf("authentication timed out")
	}
}

// RefreshToken refreshes the access token using the stored refresh token.
func (c *OAuthClient) RefreshToken() (*TokenResponse, error) {
	refreshToken, err := c.Keyring.LoadRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("no refresh token found: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	tokenURL := fmt.Sprintf("%s/api/mcp/oauth/token", c.BaseURL)
	resp, err := c.HTTPClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	// Update stored refresh token if a new one was provided
	if tokens.RefreshToken != "" {
		_ = c.Keyring.SaveRefreshToken(tokens.RefreshToken)
	}

	return &tokens, nil
}

// registerClient registers a dynamic OAuth client.
func (c *OAuthClient) registerClient(ctx context.Context) (string, error) {
	registerURL := fmt.Sprintf("%s/api/mcp/oauth/register", c.BaseURL)

	reqBody := map[string]interface{}{
		"client_name":   "MYN CLI",
		"client_uri":    "https://github.com/mindyournow/myn-cli",
		"grant_types":   []string{"authorization_code", "refresh_token"},
		"redirect_uris": []string{"http://localhost/callback"},
		"scope":         "mcp",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", registerURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("client registration failed: %s", string(body))
	}

	var reg ClientRegistration
	if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
		return "", err
	}

	return reg.ClientID, nil
}

// buildAuthURL builds the authorization URL.
func (c *OAuthClient) buildAuthURL(clientID, codeChallenge, state string, port int) string {
	authURL := fmt.Sprintf("%s/api/mcp/oauth/authorize", c.BaseURL)
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	u, _ := url.Parse(authURL)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	q.Set("scope", "mcp")
	u.RawQuery = q.Encode()

	return u.String()
}

// startCallbackServer starts a local HTTP server to receive the OAuth callback.
func (c *OAuthClient) startCallbackServer(expectedState string, codeChan chan<- string, errChan chan<- error) (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port

	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/callback" {
			http.NotFound(w, r)
			return
		}

		state := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		errorParam := r.URL.Query().Get("error")

		if errorParam != "" {
			errChan <- fmt.Errorf("OAuth error: %s", errorParam)
			http.Error(w, "Authentication failed", http.StatusBadRequest)
			return
		}

		if state != expectedState {
			errChan <- fmt.Errorf("state mismatch")
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			http.Error(w, "No code", http.StatusBadRequest)
			return
		}

		// Success response
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html><head><title>Authentication Successful</title>
<style>
body{font-family:system-ui,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#0d1117;color:#c9d1d9}
.container{text-align:center;padding:2rem;background:#161b22;border-radius:8px;border:1px solid #30363d}
h1{color:#3fb950;margin:0 0 1rem}
p{margin:0;color:#8b949e}
</style></head>
<body>
<div class="container">
<h1>✓ Authentication Successful</h1>
<p>You can close this window and return to the terminal.</p>
</div>
</body></html>`))

		codeChan <- code
	})}

	go server.Serve(listener)

	return port, nil
}

// exchangeCode exchanges the authorization code for tokens.
func (c *OAuthClient) exchangeCode(ctx context.Context, clientID, code, codeVerifier string, port int) (*TokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/api/mcp/oauth/token", c.BaseURL)
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

// generateCodeVerifier generates a random PKCE code verifier.
func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// generateCodeChallenge generates the S256 code challenge from a verifier.
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// generateState generates a random state parameter.
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
