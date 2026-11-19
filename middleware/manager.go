package middleware

import "net/http"

// Middleware defines the standard function signature for HTTP middleware.
// It accepts a handler and returns a new handler that wraps it.
type Middleware func(http.Handler) http.Handler

// Manager manages a chain of middlewares to avoid deeply nested wrapping pyramids.
type Manager struct {
	middlewares []Middleware
}

// NewManager creates a new empty Middleware Manager.
func NewManager() *Manager {
	return &Manager{
		middlewares: make([]Middleware, 0),
	}
}

// Use appends one or more middlewares to the chain.
// Middlewares are executed in the order they are passed to Use.
func (m *Manager) Use(middlewares ...Middleware) {
	m.middlewares = append(m.middlewares, middlewares...)
}

// Wrap applies all registered middlewares to the provided handler.
// It iterates in reverse so that the first middleware registered
// becomes the outermost wrapper and executes first on incoming requests.
func (m *Manager) Wrap(handler http.Handler) http.Handler {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	return handler
}
