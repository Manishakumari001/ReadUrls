package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadURL_Success(t *testing.T) {
	// Simulate a successful HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer server.Close()

	ctx := context.Background()
	data, err := downloadURL(ctx, server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("Expected 'hello world', got %s", string(data))
	}
}

func TestDownloadURL_Fail(t *testing.T) {
	ctx := context.Background()
	_, err := downloadURL(ctx, "http://invalid_url")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}
