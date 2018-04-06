package nsqd

import "sync"

type Channel struct {
	sync.RWMutex
}