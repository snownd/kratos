package mux

import "net/http"

type Muxer interface {
	Vars(r *http.Request) map[string]string
	Handle(path string, handler http.Handler)
	VerboseHandle(method, path string, handler http.Handler)
	HandleFunc(path string, f http.HandlerFunc)
	HandlePrefix(prefix string, h http.Handler)
	HandleHeader(key, val string, h http.HandlerFunc)
	WithOptions(opts *Options)
	GetPathTemplate(r *http.Request) string
	WalkRoute(fn func(path, method string) error) error
	Use(mwf ...MiddlewareFunc)
	GetHandler() http.Handler
}

type Options struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	StrictSlash             bool
	Prefix                  string
}

type MiddlewareFunc func(http.Handler) http.Handler
