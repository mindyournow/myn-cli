package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamPost_Basic(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "data: hello\n\n")
		_, _ = io.WriteString(w, "data: world\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	var received []string
	err := client.StreamPost(context.Background(), "/stream", nil, func(e SSEEvent) error {
		received = append(received, e.Data)
		return nil
	})
	if err != nil {
		t.Fatalf("StreamPost() error = %v", err)
	}
	if len(received) != 2 {
		t.Errorf("got %d events, want 2", len(received))
	}
	if received[0] != "hello" || received[1] != "world" {
		t.Errorf("got %v, want [hello world]", received)
	}
}

func TestStreamPost_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"unauthorized"}`)
	}))
	defer srv.Close()

	client := NewClient(srv.URL)
	err := client.StreamPost(context.Background(), "/stream", nil, func(e SSEEvent) error { return nil })
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
	if !strings.Contains(err.Error(), "401") {
		t.Errorf("error = %v, want 401", err)
	}
}

func TestReadSSE_MultiLineData(t *testing.T) {
	body := "data: line1\ndata: line2\n\ndata: [DONE]\n\n"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	var events []SSEEvent
	err := readSSE(context.Background(), resp, func(e SSEEvent) error {
		events = append(events, e)
		return nil
	})
	if err != nil {
		t.Fatalf("readSSE() error = %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
	if events[0].Data != "line1\nline2" {
		t.Errorf("data = %q, want %q", events[0].Data, "line1\nline2")
	}
}

func TestReadSSE_WithEventAndID(t *testing.T) {
	body := "id: 42\nevent: message\ndata: payload\n\ndata: [DONE]\n\n"
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	var events []SSEEvent
	_ = readSSE(context.Background(), resp, func(e SSEEvent) error {
		events = append(events, e)
		return nil
	})
	if len(events) != 1 || events[0].ID != "42" || events[0].Event != "message" || events[0].Data != "payload" {
		t.Errorf("got %+v", events)
	}
}
