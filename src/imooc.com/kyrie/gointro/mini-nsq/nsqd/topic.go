package nsqd

import "sync"

type Topic struct {
	sync.RWMutex
	name string
	channelMap map[string]*Channel
}