package protocol

import (
	"net"
	//"io"
)

type Protocol interface {
	IOLoop(conn net.Conn) error
}

//func SendResponse(w io.Writer, data []byte) (int, error)  {
//
//}