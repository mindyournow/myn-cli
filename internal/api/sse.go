package api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// SSEEvent represents a Server-Sent Event.
type SSEEvent struct {
	Event string
	Data  string
	ID    string
	Retry int
}

// SSEReader reads Server-Sent Events from a response body.
type SSEReader struct {
	scanner *bufio.Scanner
	reader  io.ReadCloser
}

// NewSSEReader creates a new SSE reader from an HTTP response.
func NewSSEReader(resp *http.Response) *SSEReader {
	return &SSEReader{
		scanner: bufio.NewScanner(resp.Body),
		reader:  resp.Body,
	}
}

// Close closes the underlying response body.
func (r *SSEReader) Close() error {
	return r.reader.Close()
}

// Next reads the next SSE event.
// Returns io.EOF when the stream ends.
func (r *SSEReader) Next() (*SSEEvent, error) {
	event := &SSEEvent{}

	for r.scanner.Scan() {
		line := r.scanner.Text()

		// Empty line indicates end of event
		if line == "" {
			if event.Data != "" || event.Event != "" {
				return event, nil
			}
			continue
		}

		// Parse field
		if idx := strings.IndexByte(line, ':'); idx >= 0 {
			field := line[:idx]
			value := line[idx+1:]

			// Strip leading space from value
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}

			switch field {
			case "event":
				event.Event = value
			case "data":
				if event.Data != "" {
					event.Data += "\n"
				}
				event.Data += value
			case "id":
				event.ID = value
			case "retry":
				// Ignored for now
			}
		}
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	return nil, io.EOF
}

// ReadStream reads all events from the stream and returns them as a slice.
// This is useful for non-streaming use cases where you want to collect all events.
func (r *SSEReader) ReadStream() ([]SSEEvent, error) {
	defer r.Close()

	var events []SSEEvent
	for {
		event, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return events, err
		}
		events = append(events, *event)
	}

	return events, nil
}

// StreamEvents streams SSE events to the provided channel.
// The channel is closed when the stream ends or the context is cancelled.
func (r *SSEReader) StreamEvents(ctx context.Context, ch chan<- SSEEvent) error {
	defer close(ch)
	defer r.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		event, err := r.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		ch <- *event
	}
}

// SSEToStrings converts SSE events to a slice of strings (just the data fields).
func SSEToStrings(events []SSEEvent) []string {
	var result []string
	for _, e := range events {
		result = append(result, e.Data)
	}
	return result
}

// SSEToString concatenates all event data into a single string.
func SSEToString(events []SSEEvent) string {
	var buf bytes.Buffer
	for _, e := range events {
		buf.WriteString(e.Data)
	}
	return buf.String()
}

// ValidateSSEResponse checks if a response is a valid SSE stream.
func ValidateSSEResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SSE request failed with status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		return fmt.Errorf("expected text/event-stream, got %s", contentType)
	}

	return nil
}
