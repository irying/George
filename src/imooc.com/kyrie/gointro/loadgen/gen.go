package loadgen

import (
	"imooc.com/kyrie/gointro/loadgen/lib"
	"time"
)

type myGenerator struct {
	caller lib.Caller
	timeoutNs time.Duration
	lps uint32
	durationNs time.Duration  // 负载持续时间，单位纳秒
	tickets lib.GoTickets
	stopSign chan byte
	cancelSign byte
	endSign chan byte
	callCount uint64
	status lib.GenStatus
	resultCh chan *lib.CallResult // 数组和切片不是并发安全的，要用通道
}

func NewGenerator()  {
	
}
