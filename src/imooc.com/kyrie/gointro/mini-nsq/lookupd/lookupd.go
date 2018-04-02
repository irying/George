package lookupd

import (
	"sync"
	"net"
	"utils/waitwraper"
)

type Lookupd struct {
	sync.RWMutex
	opts *Options
	tcpListener net.Listener
	httpListener net.Listener
	waitGroup waitwraper.WaitGroupWrapper
}

func New(opts *Options) *Lookupd {
	return &Lookupd{
		opts: opts,
	}
}