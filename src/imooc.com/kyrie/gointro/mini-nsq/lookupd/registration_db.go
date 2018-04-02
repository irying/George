package lookupd

import "sync"

type Registration struct {
	Category string
	Key string
	SubKey string
}

type RegistrationDB struct {
	sync.RWMutex
	registrationMap map[RegistrationDB]Producers
}

type Producer struct {
	peerInfo *PeerInfo
}

type Producers []*Producer

type PeerInfo struct {
	lastUpdate int64
	id string
	RemoteAddress string
	HostName string
	TcpPort int
	HttpPort int
}
