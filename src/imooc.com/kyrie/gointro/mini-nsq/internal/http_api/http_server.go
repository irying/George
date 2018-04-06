package http_api

import (
	"net"
	"net/http"
)

func Serve(listener net.Listener, handler http.Handler)  {
	server := &http.Server{
		Handler:handler,
	}
	err := server.Serve(listener)
	if err != nil {
		// TODO log
	}
}