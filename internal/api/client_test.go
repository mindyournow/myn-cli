package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://api.example.com")

	if client.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q, want %q", client.BaseURL, "https://api.example.com")
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient should not be nil")
	}
}

func TestClient_SetToken(t *testing.T) {
	client := NewClient("https://api.example.com")
	client.SetToken("test-token")

	if client.getToken() != "test-token" {
		t.Error("Token should be set correctly")
	}
}

func TestClient_SetToken_Concurrent(t *testing.T) {
	client := NewClient("https://api.example.com")

	// Test concurrent access to token (B41 fix - mutex protection)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client.SetToken(fmt.Sprintf("token-%d", i))
			_ = client.getToken()
		}(i)
	}
	wg.Wait()
}

func TestClient_DoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var result map[string]string
	if err := resp.DecodeJSON(&result); err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}

	if result["message"] != "success" {
		t.Errorf("message = %q, want %q", result["message"], "success")
	}
}

func TestClient_DoRequest_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	requestBody := map[string]string{"key": "value"}
	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   "/test",
		Body:   requestBody,
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	var result map[string]string
	if err := resp.DecodeJSON(&result); err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("key = %q, want %q", result["key"], "value")
	}
}

func TestClient_DoRequest_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"param1": query.Get("param1"),
			"param2": query.Get("param2"),
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
		QueryParams: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	var result map[string]string
	if err := resp.DecodeJSON(&result); err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}

	if result["param1"] != "value1" {
		t.Errorf("param1 = %q, want %q", result["param1"], "value1")
	}
	if result["param2"] != "value2" {
		t.Errorf("param2 = %q, want %q", result["param2"], "value2")
	}
}

func TestClient_DoRequest_WithAuthorization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"auth": auth})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	client.SetToken("test-token")
	ctx := context.Background()

	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	var result map[string]string
	if err := resp.DecodeJSON(&result); err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}

	if result["auth"] != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", result["auth"], "Bearer test-token")
	}
}

func TestClient_DoRequest_InvalidURL(t *testing.T) {
	client := NewClient("://invalid-url")
	ctx := context.Background()

	_, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err == nil {
		t.Error("DoRequest() should return error for invalid URL (B8 fix)")
	}
}

func TestClient_DoRequest_RetryOnServerError(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	if atomic.LoadInt32(&requestCount) < 2 {
		t.Error("Should have retried on server error")
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_DoRequest_NoRetryOnPOST(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	_, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   "/test",
		Body:   map[string]string{"key": "value"},
	})
	if err == nil {
		t.Error("DoRequest() should return error on server error")
	}

	// Should only make 1 request (no retry for POST) (B9 fix)
	if atomic.LoadInt32(&requestCount) != 1 {
		t.Errorf("Request count = %d, want 1 (POST should not retry)", atomic.LoadInt32(&requestCount))
	}
}

func TestClient_DoRequest_ClientErrorNoRetry(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	_, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err == nil {
		t.Error("DoRequest() should return error on 404")
	}

	if atomic.LoadInt32(&requestCount) != 1 {
		t.Errorf("Request count = %d, want 1 (client errors should not retry)", atomic.LoadInt32(&requestCount))
	}
}

func TestClient_DoRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err == nil {
		t.Error("DoRequest() should return error on context timeout")
	}
}

func TestClient_DoRequest_RateLimitRetry(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	resp, err := client.DoRequest(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err != nil {
		t.Fatalf("DoRequest() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	_, err := client.Get(ctx, "/test", nil)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	ctx := context.Background()

	_, err := client.Post(ctx, "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}
}

func TestResponse_IsSuccess(t *testing.T) {
	tests := []struct {
		statusCode int
		want       bool
	}{
		{200, true},
		{201, true},
		{299, true},
		{199, false},
		{300, false},
		{400, false},
		{500, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			r := &Response{StatusCode: tt.statusCode}
			if got := r.IsSuccess(); got != tt.want {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_DecodeJSON(t *testing.T) {
	r := &Response{
		Body: []byte(`{"key": "value"}`),
	}

	var result map[string]string
	if err := r.DecodeJSON(&result); err != nil {
		t.Fatalf("DecodeJSON() error = %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("key = %q, want %q", result["key"], "value")
	}
}

func TestResponse_IsNotFound(t *testing.T) {
	tests := []struct {
		statusCode int
		want       bool
	}{
		{404, true},
		{200, false},
		{500, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			r := &Response{StatusCode: tt.statusCode}
			if got := r.IsNotFound(); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_ParseError(t *testing.T) {
	r := &Response{
		StatusCode: 400,
		Body:       []byte(`{"error": "bad_request", "message": "Invalid input"}`),
	}

	errResp := r.ParseError()
	if errResp.Error != "bad_request" {
		t.Errorf("Error = %q, want %q", errResp.Error, "bad_request")
	}
	if errResp.Message != "Invalid input" {
		t.Errorf("Message = %q, want %q", errResp.Message, "Invalid input")
	}
	if errResp.Code != 400 {
		t.Errorf("Code = %d, want %d", errResp.Code, 400)
	}
}

func TestResponse_ParseError_NonJSON(t *testing.T) {
	r := &Response{
		StatusCode: 500,
		Body:       []byte("Internal Server Error"),
	}

	errResp := r.ParseError()
	if errResp.Error != "unknown_error" {
		t.Errorf("Error = %q, want %q", errResp.Error, "unknown_error")
	}
	if !strings.Contains(errResp.Message, "Internal Server Error") {
		t.Errorf("Message = %q, should contain 'Internal Server Error'", errResp.Message)
	}
}

func TestResponse_ParseError_DefaultCode(t *testing.T) {
	r := &Response{
		StatusCode: 500,
		Body:       []byte(`{"error": "error_without_code"}`),
	}

	errResp := r.ParseError()
	if errResp.Code != 500 {
		t.Errorf("Code = %d, want %d (should default to StatusCode)", errResp.Code, 500)
	}
}
