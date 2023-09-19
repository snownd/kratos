// go:build !bunmux
package mux

import (
	"net/http"

	"github.com/gorilla/mux"
)

var _ Muxer = (*gorillaMux)(nil)

func NewMuxer() Muxer {
	return &gorillaMux{mux.NewRouter()}
}

type gorillaMux struct {
	router *mux.Router
}

func (r *gorillaMux) Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}

func (r *gorillaMux) GetPathTemplate(req *http.Request) string {
	pathTemplate := req.URL.Path
	if route := mux.CurrentRoute(req); route != nil {
		// /path/123 -> /path/{id}
		v, err := route.GetPathTemplate()
		if err == nil {
			pathTemplate = v
		}
	}
	return pathTemplate
}

func (r *gorillaMux) Handle(path string, handler http.Handler) {
	r.router.Handle(path, handler)
}

func (r *gorillaMux) VerboseHandle(method, path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(method)
}

func (r *gorillaMux) HandleFunc(path string, f http.HandlerFunc) {
	r.router.HandleFunc(path, f)
}

func (r *gorillaMux) HandlePrefix(prefix string, h http.Handler) {
	r.router.PathPrefix(prefix).Handler(h)
}

func (r *gorillaMux) HandleHeader(key string, val string, h http.HandlerFunc) {
	r.router.Headers(key, val).Handler(h)
}

func (r *gorillaMux) WithOptions(opts *Options) {
	if opts.Prefix != "" {
		r.router = r.router.PathPrefix(opts.Prefix).Subrouter()
	}
	if opts.MethodNotAllowedHandler != nil {
		r.router.MethodNotAllowedHandler = opts.MethodNotAllowedHandler
	}
	if opts.NotFoundHandler != nil {
		r.router.NotFoundHandler = opts.NotFoundHandler
	}
	if opts.StrictSlash {
		r.router.StrictSlash(true)
	}
}

func (r *gorillaMux) WalkRoute(fn func(path, method string) error) error {
	return r.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, err := route.GetMethods()
		if err != nil {
			return nil // ignore no methods
		}
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		for _, method := range methods {
			if err := fn(path, method); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gorillaMux) Use(mwf ...MiddlewareFunc) {
	gmwf := make([]mux.MiddlewareFunc, 0)
	for _, f := range mwf {
		gmwf = append(gmwf, mux.MiddlewareFunc(f))
	}
	r.router.Use(gmwf...)
}

func (r *gorillaMux) GetHandler() http.Handler {
	return r.router
}
