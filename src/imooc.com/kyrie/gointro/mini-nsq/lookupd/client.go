package lookupd

import "net"

type Client struct {
	net.Conn
	peerInfo *PeerInfo
}

func NewClient(conn net.Conn) *Client  {
	return &Client{
		Conn:conn,
	}
}
