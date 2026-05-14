package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(2080)
	if s == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestAddRoute(t *testing.T) {
	s := New(2080)
	s.AddRoute("a.example.com", "localhost:3000")

	s.mu.RLock()
	target, ok := s.routes["a.example.com"]
	s.mu.RUnlock()

	if !ok {
		t.Fatal("route not found after AddRoute")
	}
	if target != "localhost:3000" {
		t.Errorf("target: got %q, want %q", target, "localhost:3000")
	}
}

func TestRemoveRoute(t *testing.T) {
	s := New(2080)
	s.AddRoute("a.example.com", "localhost:3000")
	s.RemoveRoute("a.example.com")

	s.mu.RLock()
	_, ok := s.routes["a.example.com"]
	s.mu.RUnlock()

	if ok {
		t.Error("route still present after RemoveRoute")
	}
}

func TestServeHTTP_proxiesRequest(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "hello from backend")
	}))
	defer backend.Close()

	s := New(2080)
	s.AddRoute("a.example.com", strings.TrimPrefix(backend.URL, "http://"))

	req := httptest.NewRequest(http.MethodGet, "http://a.example.com/path", nil)
	req.Host = "a.example.com"
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if body := w.Body.String(); !strings.Contains(body, "hello from backend") {
		t.Errorf("body: got %q, want 'hello from backend'", body)
	}
}

func TestServeHTTP_unknownRoute(t *testing.T) {
	s := New(2080)

	req := httptest.NewRequest(http.MethodGet, "http://unknown.example.com/", nil)
	req.Host = "unknown.example.com"
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type: got %q, want text/html", ct)
	}
	if !strings.Contains(w.Body.String(), "Tunnel Not Found") {
		t.Error("body does not contain custom error message")
	}
}

func TestServeHTTP_unreachableUpstream(t *testing.T) {
	s := New(2080)
	s.AddRoute("a.example.com", "localhost:19999") // nothing listening

	req := httptest.NewRequest(http.MethodGet, "http://a.example.com/", nil)
	req.Host = "a.example.com"
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status: got %d, want 502", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Service Unreachable") {
		t.Error("body does not contain custom 502 error message")
	}
}

func TestServeHTTP_hostWithPort(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	s := New(2080)
	s.AddRoute("a.example.com", strings.TrimPrefix(backend.URL, "http://"))

	req := httptest.NewRequest(http.MethodGet, "http://a.example.com:2080/", nil)
	req.Host = "a.example.com:2080"
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("host with port not routed correctly: got %d", w.Code)
	}
}

func TestStartStop(t *testing.T) {
	port := freePort(t)
	s := New(port)

	errCh := make(chan error, 1)
	go func() { errCh <- s.Start() }()

	// wait for the server to be ready
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 50*time.Millisecond)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		t.Fatalf("request to running server: %v", err)
	}
	resp.Body.Close()

	if err := s.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Start returned error after Stop: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("Start did not return after Stop")
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("finding free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}
