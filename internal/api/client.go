package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client handles HTTP communication with the MYN backend.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	apiKey     string
	token      string
	retries    int
}

// RequestOptions contains optional parameters for requests.
type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
}

// NewClient creates a new API client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retries: 3,
	}
}

// SetAPIKey sets the API key for authentication.
func (c *Client) SetAPIKey(key string) {
	c.apiKey = key
}

// SetToken sets the Bearer token for authenticated requests.
func (c *Client) SetToken(token string) {
	c.token = token
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, opts *RequestOptions) (*http.Response, error) {
	return c.doRequest(ctx, "GET", path, nil, opts)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.doRequest(ctx, "POST", path, body, opts)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.doRequest(ctx, "PUT", path, body, opts)
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.doRequest(ctx, "PATCH", path, body, opts)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, opts *RequestOptions) (*http.Response, error) {
	return c.doRequest(ctx, "DELETE", path, nil, opts)
}

// doRequest performs an HTTP request with retries and auth injection.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			delay := time.Duration(1<<(attempt-1)) * time.Second
			time.Sleep(delay)
		}

		resp, err := c.executeRequest(ctx, method, path, body, opts)
		if err != nil {
			lastErr = err
			continue
		}

		// Check if we should retry
		if c.shouldRetry(resp) {
			lastErr = fmt.Errorf("server returned %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retries+1, lastErr)
}

// executeRequest executes a single HTTP request.
func (c *Client) executeRequest(ctx context.Context, method, path string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	url := c.buildURL(path, opts)

	var bodyReader io.Reader
	if body != nil {
		if r, ok := body.(io.Reader); ok {
			bodyReader = r
		} else {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("marshaling request body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBody)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Set default headers
	req.Header.Set("Accept", "application/json")
	if body != nil {
		if _, ok := body.(io.Reader); !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Add auth header
	c.addAuthHeader(req)

	// Add custom headers
	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	return c.HTTPClient.Do(req)
}

// buildURL constructs the full URL with query parameters.
func (c *Client) buildURL(path string, opts *RequestOptions) string {
	u := c.BaseURL + path

	if opts != nil && len(opts.Query) > 0 {
		parsedURL, _ := url.Parse(u)
		q := parsedURL.Query()
		for k, v := range opts.Query {
			q.Set(k, v)
		}
		parsedURL.RawQuery = q.Encode()
		u = parsedURL.String()
	}

	return u
}

// addAuthHeader adds the appropriate authentication header.
func (c *Client) addAuthHeader(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	} else if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

// shouldRetry determines if a request should be retried based on the response.
func (c *Client) shouldRetry(resp *http.Response) bool {
	// Retry on server errors and rate limiting
	return resp.StatusCode >= 500 || resp.StatusCode == 429
}

// DecodeJSON decodes a JSON response body into the target.
func DecodeJSON(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

// ReadBody reads the response body as a string.
func ReadBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}
