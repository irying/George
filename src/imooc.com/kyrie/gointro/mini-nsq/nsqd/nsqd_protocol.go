package nsqd

import (
	"net"
	"sync/atomic"
)

type Protocol struct {
	ctx *context
} 

func (p *Protocol)IOLoop(conn net.Conn) error {
	var err error
	var line []byte
	clientID := atomic.AddInt64(&p.ctx.nsqd.clientIDSequence, 1)
	for {
		line, err :=
	}
}