package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// SSEEvent represents a single server-sent event.
type SSEEvent struct {
	ID    string
	Event string
	Data  string
}

// SSEHandler is called for each SSE event. Return an error to stop streaming.
type SSEHandler func(event SSEEvent) error

// StreamPost performs a POST request and reads the response as a server-sent event stream.
// The handler is called for each event. Streaming ends when [DONE] is received or context is cancelled.
func (c *Client) StreamPost(ctx context.Context, path string, body interface{}, handler SSEHandler) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	fullURL, err := buildURL(c.BaseURL, path)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if apiKey := c.getAPIKey(); apiKey != "" {
		req.Header.Set("X-API-KEY", apiKey)
	} else if token := c.getToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("request failed: %s - %s", resp.Status, string(errBody))
	}

	return readSSE(ctx, resp, handler)
}

// readSSE parses text/event-stream from an HTTP response, calling handler for each event.
func readSSE(ctx context.Context, resp *http.Response, handler SSEHandler) error {
	scanner := bufio.NewScanner(resp.Body)
	var event SSEEvent
	gotDone := false

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()

		if line == "" {
			// Empty line = dispatch event
			if event.Data != "" {
				if event.Data == "[DONE]" {
					gotDone = true
					return nil
				}
				if err := handler(event); err != nil {
					return err
				}
			}
			event = SSEEvent{}
			continue
		}

		if strings.HasPrefix(line, "id:") {
			event.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		} else if strings.HasPrefix(line, "event:") {
			event.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if event.Data == "" {
				event.Data = data
			} else {
				event.Data += "\n" + data
			}
		}
		// Skip comment lines (starting with ':') and unknown fields
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading SSE stream: %w", err)
	}
	if !gotDone {
		return fmt.Errorf("SSE stream ended without [DONE] signal (incomplete response)")
	}
	return nil
}

// buildURL constructs a full URL from base and path.
func buildURL(base, path string) (string, error) {
	full, err := url.JoinPath(base, path)
	if err != nil {
		return "", fmt.Errorf("failed to build URL: %w", err)
	}
	return full, nil
}
