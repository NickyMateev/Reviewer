package web

import "net/http"

// HandleFunc defines a function type needed for Route constructions
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// ServeHTTP allows normal HandlerFunc functions to behave like Handler implementations
func (hf HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hf(w, r)
}

// Route defines a route that will be attached to a router
type Route struct {
	Path       string
	Method     string
	HandleFunc HandleFunc
}

// Controller defines a set of HTTP routes
//go:generate counterfeiter . Controller
type Controller interface {
	Routes() []Route
}

// API defines a set of controllers which constitute the whole API of the application
//go:generate counterfeiter . API
type API interface {
	Controllers() []Controller
}
