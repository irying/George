package loadgen

import (
	"imooc.com/kyrie/gointro/loadgen/lib"
	"time"
	"log"
	"fmt"
	"errors"
	"math"
)

// 3类东西组成，
// 调用前的调用者，调用时间，go池数量，调用器状态
// 调用中需要统计的lps,持续时间，调用停止等信号
// 调用结果
type myGenerator struct {
	caller     lib.Caller
	timeoutNs  time.Duration
	lps        uint32
	durationNs time.Duration        // 负载持续时间，单位纳秒
	concurrency uint32
	tickets    lib.GoTickets
	stopSign   chan byte            // 停止信号的传递通道
	cancelSign byte                 // 取消发送后续结果的信号
	endSign    chan byte            // 完结信号的传递通道，同时被用于传递调用执行计数。
	callCount  uint64
	status     lib.GenStatus
	resultCh   chan *lib.CallResult // 数组和切片不是并发安全的，要用通道
}

func NewGenerator(
caller lib.Caller,
timeoutNs time.Duration,
lps uint32,
durationNs time.Duration,
resultCh chan *lib.CallResult) (lib.Generator, error) {
	log.Println("New a load generator...")
	log.Println("Checking the parameters...")
	// Checking the parameters
	var errMsg string
	if caller == nil {
		errMsg = fmt.Sprintln("Invalid caller!")
	}

	if timeoutNs == 0 {
		errMsg
	}
	if timeoutNs == 0 {
		errMsg = fmt.Sprintln("Invalid timeoutNs!")
	}
	if lps == 0 {
		errMsg = fmt.Sprintln("Invalid lps(load per second)!")
	}
	if durationNs == 0 {
		errMsg = fmt.Sprintln("Invalid durationNs!")
	}
	if resultCh == nil {
		errMsg = fmt.Sprintln("Invalid result channel!")
	}
	if errMsg != "" {
		return nil, errors.New(errMsg)
	}
	gen := &myGenerator{
		caller: caller,
		timeoutNs:timeoutNs,
		lps:lps,
		durationNs:durationNs,
		stopSign:make(chan byte, 1),
		cancelSign:0,
		status:lib.STATUS_ORIGINAL,
		resultCh:resultCh,
	}
	log.Printf("Passed. (timeoutNs=%v, lps=%d, durationNs=%v)", timeoutNs, lps, durationNs)
	err := gen.init()
	if err != nil {
		return nil, err
	}

	return gen, nil
}

func (gen *myGenerator) init() error {
	log.Println("Initializing the load generator...")
	// 载荷的并发量 = 载荷的响应超时时间 ／ 载荷的发送间隔时间
	// TODO: uint64 uint32
	var total64 int64 = int64(gen.timeoutNs)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets
	log.Printf("Initialized. (concurrency=%d)", gen.concurrency)
	return nil
}
