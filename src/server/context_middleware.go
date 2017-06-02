package main

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/pborman/uuid"
)

// A ContextMiddleware ensure each request is given an initialized context.
type ContextMiddleware struct {
	ctx Context
}

func NewContextMiddleware(ctx Context) *ContextMiddleware {
	return &ContextMiddleware{
		ctx: ctx,
	}
}

func (m *ContextMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, n http.HandlerFunc) {
	// Generate an ID for the request.
	requestID := uuid.New()

	// Generate a custom context for the request.
	ctx := m.ctx
	ctx.RequestID = requestID
	ctx.Logger = log.With(ctx.Logger, "request", requestID)

	// Generate a new context for the request.
	clock.Lock()
	cmap[r] = &ctx
	clock.Unlock()

	// Execute the next handlers.
	n(w, r)

	// Remove the context from the map.
	clock.Lock()
	delete(cmap, r)
	clock.Unlock()
}
