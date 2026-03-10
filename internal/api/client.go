package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client handles HTTP communication with the MYN backend.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	mu         sync.RWMutex // protects token field (B41 fix - mutex protection)
	token      string
}

// NewClient creates a new API client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetToken sets the Bearer token for authenticated requests.
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

// getToken safely retrieves the current token.
func (c *Client) getToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// RequestOptions contains options for making HTTP requests.
type RequestOptions struct {
	Method      string
	Path        string
	Body        interface{}
	QueryParams map[string]string
	Headers     map[string]string
}

// Response wraps the HTTP response with additional metadata.
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// DoRequest performs an HTTP request with retry logic and proper error handling.
// Callers should NOT close resp.Body as this method handles it internally.
func (c *Client) DoRequest(ctx context.Context, opts RequestOptions) (*Response, error) {
	// Validate and construct URL (B8 fix - proper URL parsing error handling)
	fullURL, err := url.JoinPath(c.BaseURL, opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add query parameters
	if len(opts.QueryParams) > 0 {
		q := u.Query()
		for key, value := range opts.QueryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Serialize body if provided
	var bodyReader io.Reader
	var bodyBytes []byte
	if opts.Body != nil {
		if bodyStr, ok := opts.Body.(string); ok {
			bodyBytes = []byte(bodyStr)
		} else {
			bodyBytes, err = json.Marshal(opts.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Determine if request is idempotent (safe for retry)
	isIdempotent := opts.Method == http.MethodGet ||
		opts.Method == http.MethodHead ||
		opts.Method == http.MethodOptions ||
		opts.Method == http.MethodTrace

	// Perform request with retry logic
	var lastErr error
	maxRetries := 3
	baseDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// For non-idempotent requests (POST, PUT, PATCH, DELETE), don't retry
			// because the body would be consumed (B9 fix)
			if !isIdempotent {
				break
			}

			// Exponential backoff with context cancellation check (warning fix)
			delay := baseDelay * time.Duration(1<<uint(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
				// Continue to retry
			}
		}

		// Reset body reader for retry (only safe for idempotent requests)
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, opts.Method, u.String(), bodyReader)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set default headers
		req.Header.Set("Accept", "application/json")
		if opts.Body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// Set authorization header if token is available
		if token := c.getToken(); token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		// Set custom headers
		for key, value := range opts.Headers {
			req.Header.Set(key, value)
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			// Retry on network errors
			continue
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// Check for success
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return &Response{
				StatusCode: resp.StatusCode,
				Headers:    resp.Header,
				Body:       body,
			}, nil
		}

		// Handle rate limiting (429) - parse Retry-After header (warning fix)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				// Parse Retry-After as seconds
				var seconds int
				if _, err := fmt.Sscanf(retryAfter, "%d", &seconds); err == nil {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(time.Duration(seconds) * time.Second):
						continue
					}
				}
			}
		}

		// Don't retry client errors (4xx)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, fmt.Errorf("request failed: %s - %s", resp.Status, string(body))
		}

		// Server error (5xx) - retry
		lastErr = fmt.Errorf("server error: %s - %s", resp.Status, string(body))
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries, lastErr)
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string) (*Response, error) {
	return c.DoRequest(ctx, RequestOptions{
		Method:      http.MethodGet,
		Path:        path,
		QueryParams: queryParams,
	})
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.DoRequest(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.DoRequest(ctx, RequestOptions{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.DoRequest(ctx, RequestOptions{
		Method: http.MethodPatch,
		Path:   path,
		Body:   body,
	})
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.DoRequest(ctx, RequestOptions{
		Method: http.MethodDelete,
		Path:   path,
	})
}

// IsSuccess returns true if the response indicates success (2xx status).
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// UnmarshalJSON unmarshals the response body as JSON into the provided target.
func (r *Response) UnmarshalJSON(target interface{}) error {
	return json.Unmarshal(r.Body, target)
}

// IsNotFound returns true if the response is a 404.
func (r *Response) IsNotFound() bool {
	return r.StatusCode == http.StatusNotFound
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ParseError attempts to parse the response body as an ErrorResponse.
func (r *Response) ParseError() *ErrorResponse {
	var errResp ErrorResponse
	if err := json.Unmarshal(r.Body, &errResp); err != nil {
		// If we can't parse as JSON, use the raw body as the message
		return &ErrorResponse{
			Error:   "unknown_error",
			Message: strings.TrimSpace(string(r.Body)),
			Code:    r.StatusCode,
		}
	}
	if errResp.Code == 0 {
		errResp.Code = r.StatusCode
	}
	return &errResp
}
