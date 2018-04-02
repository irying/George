package protocol

import "net"

type TcpHandler interface {
	Handler(net.Conn)
}