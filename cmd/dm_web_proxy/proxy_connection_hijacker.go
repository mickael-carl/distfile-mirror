package main

import (
	"net/http"
)

// proxyConnectionHijacker implements a http.Handler that filters out
// HTTP CONNECT requests and extracts the associated TCP connections to
// a ProxyConnectionHandler. Plain non-CONNECT HTTP requests are
// forwarded to another http.Handler.
type proxyConnectionHijacker struct {
	connectionHandler ProxyConnectionHandler
	requestHandler    http.Handler
}

func NewProxyConnectionHijacker(connectionHandler ProxyConnectionHandler, requestHandler http.Handler) http.Handler {
	return &proxyConnectionHijacker{
		connectionHandler: connectionHandler,
		requestHandler:    requestHandler,
	}
}

func (sch *proxyConnectionHijacker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		w.WriteHeader(http.StatusOK)
		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported for connection", http.StatusInternalServerError)
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
		sch.connectionHandler.Handle(conn, r)
	} else {
		sch.requestHandler.ServeHTTP(w, r)
	}
}
