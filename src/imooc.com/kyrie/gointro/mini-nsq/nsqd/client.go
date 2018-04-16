package nsqd

import "net"

type client struct {
	net.Conn
	ID int64
	lenBuf [4]byte
	lenSlice []byte
	ctx *context
}

func newClient(id int64, conn net.Conn, ctx *context) *client{
	client := &client{
		ID:id,
		Conn:conn,
		ctx:ctx,
	}
	client.lenSlice = client.lenBuf[:]
	return client
}