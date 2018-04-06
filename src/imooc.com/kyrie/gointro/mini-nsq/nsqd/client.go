package nsqd

import "net"

type client struct {
	net.Conn
	ID int64
	lenBuf [4]byte
	lenSlice []byte
}

func newClient(id int64, conn net.Conn) *client{
	client := &client{
		ID:id,
		Conn:conn,
	}
	client.lenSlice = client.lenBuf[:]
	return client
}