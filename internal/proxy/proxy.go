// Package proxy implements the internal HTTP reverse proxy.
package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

// Server is the internal reverse proxy that routes by Host header.
// Routes are stored in memory; add/remove is instantaneous.
type Server struct {
	mu     sync.RWMutex
	routes map[string]string // hostname → "localhost:PORT"
	port   int
	srv    *http.Server
}

// New creates a Server that will listen on the given port.
func New(port int) *Server {
	return &Server{
		routes: make(map[string]string),
		port:   port,
	}
}

// ServeHTTP routes the request by Host header to the registered target.
// Unknown hosts return a custom 404 page; unreachable upstreams return a custom 502.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host // Host has no port component
	}
	if host == "" {
		serveErrorPage(w, http.StatusBadRequest, "Missing Host", "No Host header was present in the request.")
		return
	}

	s.mu.RLock()
	target, ok := s.routes[host]
	s.mu.RUnlock()

	if !ok {
		serveErrorPage(w, http.StatusNotFound, "Tunnel Not Found",
			"No xpoz tunnel is registered for this hostname.")
		return
	}

	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
			req.Host = target
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			serveErrorPage(w, http.StatusBadGateway, "Service Unreachable",
				"Your local service is not responding. Make sure it is running on the expected port.")
		},
	}
	rp.ServeHTTP(w, r)
}

// Start begins serving HTTP on the configured port.
// It blocks until Stop is called or a listen error occurs.
func (s *Server) Start() error {
	s.mu.Lock()
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s,
	}
	srv := s.srv
	s.mu.Unlock()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("proxy: listen on :%d: %w", s.port, err)
	}
	return nil
}

// Stop gracefully shuts down the proxy, waiting up to 5 seconds for in-flight requests.
func (s *Server) Stop() error {
	s.mu.RLock()
	srv := s.srv
	s.mu.RUnlock()

	if srv == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("proxy: shutdown: %w", err)
	}
	return nil
}

// AddRoute registers hostname → target routing (e.g. "a3f9c.example.com" → "localhost:3000").
func (s *Server) AddRoute(hostname, target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[hostname] = target
}

// RemoveRoute removes the route for hostname.
func (s *Server) RemoveRoute(hostname string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.routes, hostname)
}
