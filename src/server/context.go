package main

import (
	"net/http"
	"sync"

	"github.com/go-kit/kit/log"
)

// context enable request-associated data and structures to cross middelwares
// and handlers boundaries.
type Context struct {
	RequestID string
	Logger    log.Logger
}

var cmap = map[*http.Request]*Context{}
var clock sync.RWMutex

// Context returns the context associated with the given request, or nil if
// none.
func Ctx(r *http.Request) *Context {
	clock.RLock()
	defer clock.RUnlock()
	return cmap[r]
}
