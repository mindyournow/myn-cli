package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	// PKCE constants
	codeVerifierLength = 128
	stateLength        = 32

	// OAuth endpoints (MYN backend paths)
	registerPath = "/api/mcp/oauth/register"
	authPath     = "/api/mcp/oauth/authorize"
	tokenPath    = "/api/mcp/oauth/token"
)

// callbackResult holds the result of the OAuth callback.
type callbackResult struct {
	tokens *TokenResponse
	err    error
}

// OAuthClient handles the OAuth 2.0 flow with PKCE.
type OAuthClient struct {
	BaseURL        string
	HTTPClient     *http.Client
	TokenStore     TokenStore
	ClientID       string
	ClientSecret   string // For confidential clients, empty for public clients
	callbackResult chan callbackResult
}

// NewOAuthClient creates a new OAuth client.
func NewOAuthClient(baseURL string, tokenStore TokenStore) *OAuthClient {
	return &OAuthClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		TokenStore:     tokenStore,
		callbackResult: make(chan callbackResult, 1),
	}
}

// TokenResponse represents the OAuth token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// Authenticate performs the full OAuth flow with PKCE.
// It starts a local callback server, opens the browser, and exchanges the code for tokens.
func (c *OAuthClient) Authenticate(ctx context.Context) (*TokenResponse, error) {
	// Generate PKCE parameters using crypto/rand (not math/rand)
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	state, err := generateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Get a listener for the callback server (we need this early for registration)
	// Bind explicitly to 127.0.0.1 (not "localhost") to avoid DNS-resolution ambiguity (MED-2 fix)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	// Register the client if not already registered (passing the redirect URI)
	if c.ClientID == "" {
		if err := c.registerClient(ctx, redirectURI); err != nil {
			listener.Close()
			return nil, fmt.Errorf("failed to register client: %w", err)
		}
	}

	// Start callback server using the SAME listener (B4 fix)
	server, err := c.startCallbackServer(listener, state, codeVerifier, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	// Ensure server is shut down when we're done (B5 fix)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	// Build and open authorization URL
	authURL, err := c.buildAuthURL(codeVerifier, state, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to build auth URL: %w", err)
	}

	// Try to open the browser automatically
	opened := openBrowser(authURL)
	if opened {
		fmt.Println("Browser opened for authentication.")
	} else {
		fmt.Println("Could not open browser automatically.")
	}
	fmt.Printf("\nIf your browser didn't open, copy and paste this URL:\n\n  %s\n\n", authURL)
	fmt.Println("Waiting for authentication...")

	// Wait for callback or context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-c.callbackResult:
		if result.err != nil {
			return nil, result.err
		}
		return result.tokens, nil
	}
}

// RefreshToken exchanges a refresh token for new access tokens.
func (c *OAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.ClientID)
	if c.ClientSecret != "" {
		data.Set("client_secret", c.ClientSecret)
	}

	tokenURL, err := url.JoinPath(c.BaseURL, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build token URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("token refresh failed: %s - %s", resp.Status, string(body))
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Save the new refresh token - DON'T ignore errors (B10 fix)
	if tokens.RefreshToken != "" && c.TokenStore != nil {
		if err := c.TokenStore.SaveRefreshToken(tokens.RefreshToken); err != nil {
			return nil, fmt.Errorf("failed to save refresh token: %w", err)
		}
	}

	return &tokens, nil
}

// registerClient dynamically registers this CLI as an OAuth client.
// The redirectURI is now passed in to ensure consistency (B4 fix).
func (c *OAuthClient) registerClient(ctx context.Context, redirectURI string) error {
	registrationData := map[string]interface{}{
		"client_name":   "MYN CLI",
		"client_uri":    "https://github.com/mindyournow/myn-cli",
		"redirect_uris": []string{redirectURI},
		"grant_types":   []string{"authorization_code", "refresh_token"},
		"token_endpoint_auth_method": "none", // Public client
	}

	body, err := json.Marshal(registrationData)
	if err != nil {
		return fmt.Errorf("failed to marshal registration data: %w", err)
	}

	registerURL, err := url.JoinPath(c.BaseURL, registerPath)
	if err != nil {
		return fmt.Errorf("failed to build register URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registerURL, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return fmt.Errorf("client registration failed: %s - %s", resp.Status, string(respBody))
	}

	var result struct {
		ClientID string `json:"client_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode registration response: %w", err)
	}

	c.ClientID = result.ClientID

	// Persist client_id so subsequent processes can use it for token refresh
	if c.TokenStore != nil {
		if saver, ok := c.TokenStore.(interface{ SaveClientID(string) error }); ok {
			if err := saver.SaveClientID(result.ClientID); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to save client ID: %v\n", err)
			}
		}
	}
	return nil
}

// buildAuthURL constructs the authorization URL with PKCE parameters.
// Now returns an error instead of silently discarding them (B4/B8 fix).
func (c *OAuthClient) buildAuthURL(codeVerifier, state, redirectURI string) (string, error) {
	codeChallenge := generateCodeChallenge(codeVerifier)

	authURL, err := url.JoinPath(c.BaseURL, authPath)
	if err != nil {
		return "", fmt.Errorf("failed to build auth path: %w", err)
	}

	u, err := url.Parse(authURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse auth URL: %w", err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.ClientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	q.Set("state", state)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// startCallbackServer starts an HTTP server to receive the OAuth callback.
// Accepts a pre-created listener to ensure port consistency (B4 fix).
// Returns the server instance (for shutdown) and any error.
func (c *OAuthClient) startCallbackServer(listener net.Listener, state, codeVerifier, redirectURI string) (*http.Server, error) {
	mux := http.NewServeMux()
	// Add timeouts to prevent slow-loris / hanging goroutines (MED-1 fix)
	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		queryState := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		errorParam := r.URL.Query().Get("error")

		if errorParam != "" {
			// Sanitize errorParam before using in error message to prevent ANSI injection (MED-3 fix)
			c.callbackResult <- callbackResult{err: fmt.Errorf("oauth error: %s", sanitizeParam(errorParam))}
			http.Error(w, "Authentication failed", http.StatusBadRequest)
			return
		}

		if queryState != state {
			c.callbackResult <- callbackResult{err: fmt.Errorf("state mismatch: possible CSRF attack")}
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		tokens, err := c.exchangeCode(r.Context(), code, codeVerifier, redirectURI)
		if err != nil {
			c.callbackResult <- callbackResult{err: err}
			http.Error(w, "Token exchange failed", http.StatusInternalServerError)
			return
		}

		c.callbackResult <- callbackResult{tokens: tokens}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html><body>
<h1>Authentication Successful</h1>
<p>You can close this window and return to the CLI.</p>
</body></html>`))
	})

	// Start server in goroutine
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			// Log error but don't block
			fmt.Fprintf(os.Stderr, "Callback server error: %v\n", err)
		}
	}()

	return server, nil
}

// exchangeCode exchanges the authorization code for tokens.
func (c *OAuthClient) exchangeCode(ctx context.Context, code, codeVerifier, redirectURI string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", codeVerifier)
	data.Set("client_id", c.ClientID) // Required for public clients (B3 fix)

	tokenURL, err := url.JoinPath(c.BaseURL, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build token URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
	}

	var tokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Save refresh token - DON'T ignore errors (B10 fix)
	if tokens.RefreshToken != "" && c.TokenStore != nil {
		if err := c.TokenStore.SaveRefreshToken(tokens.RefreshToken); err != nil {
			return nil, fmt.Errorf("failed to save refresh token: %w", err)
		}
	}

	return &tokens, nil
}

// generateCodeVerifier generates a cryptographically secure random code verifier.
func generateCodeVerifier() (string, error) {
	bytes := make([]byte, codeVerifierLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// generateState generates a cryptographically secure random state parameter.
func generateState() (string, error) {
	bytes := make([]byte, stateLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// sanitizeParam strips non-printable characters from an OAuth callback parameter
// to prevent terminal escape sequence injection (MED-3 fix).
func sanitizeParam(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c < 0x7f {
			out = append(out, c)
		}
	}
	return string(out)
}

// generateCodeChallenge generates the S256 code challenge from the verifier.
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// openBrowser tries to open a URL in the user's default browser.
// Returns true if the command was launched successfully.
func openBrowser(url string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return false
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start() == nil
}
