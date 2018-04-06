package nsqd

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type httpServer struct {
	ctx *context
	router http.Handler
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	s := &httpServer{
		ctx:ctx,
		router:router,
	}

	return s
}