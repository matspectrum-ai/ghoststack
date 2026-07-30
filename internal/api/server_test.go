package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ghoststack/ghoststack/internal/runtime"
)

func TestStatusEndpoint(t *testing.T) {
	daemon := runtime.NewDaemon("", nil)
	srv := NewServer(daemon)

	server, err := srv.Start(context.Background(), "127.0.0.1:18080")
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	defer server.Close()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:18080/api/status")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
