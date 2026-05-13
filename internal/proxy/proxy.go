// Package proxy implements the internal HTTP reverse proxy.
package proxy

import (
	"fmt"
	"net/http"
	"sync"
)

// Server is the internal reverse proxy that routes by Host header.
// Routes are stored in memory; add/remove is instantaneous.
type Server struct {
	mu     sync.RWMutex
	routes map[string]string // hostname → "localhost:PORT"
	port   int
	server *http.Server
}

// New creates a Server that will listen on the given port.
func New(port int) *Server {
	return &Server{
		routes: make(map[string]string),
		port:   port,
	}
}

// Start begins serving HTTP on the configured port.
// It blocks until Stop is called.
func (s *Server) Start() error {
	// TODO: implement in Phase 1
	return fmt.Errorf("proxy: not implemented")
}

// Stop gracefully shuts down the proxy server.
func (s *Server) Stop() error {
	// TODO: implement in Phase 1
	return fmt.Errorf("proxy: not implemented")
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
