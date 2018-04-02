package lookupd

import (
	"net"
	"imooc.com/kyrie/gointro/mini-nsq/internal/protocol"
)

type tcpServer struct {
	ctx *Context
}

// TODO 打log
func (p *tcpServer) Handle(clientConn net.Conn)  {
	var prot protocol.Protocol
	prot = &LookupProtocol{ctx:p.ctx}
	err := prot.IOLoop(clientConn)
	if err != nil {
		// log
		return
	}
}