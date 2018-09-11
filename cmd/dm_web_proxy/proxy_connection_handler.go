package main

import (
	"net"
	"net/http"
)

// ProxyConnectionHandler is used by ProxyConnectionHijacker to push out
// TCP connections on which HTTP connect requests are observed.
type ProxyConnectionHandler interface {
	Handle(c net.Conn, r *http.Request)
}
