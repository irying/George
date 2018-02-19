package loadgen

import (
	"imooc.com/kyrie/gointro/loadgen/lib"
	"time"
	"log"
	"fmt"
	"errors"
	"math"
	"bytes"
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

func (gen *myGenerator) genLoad(throttle <-chan time.Time)  {
	callCount := uint64(0)
	Loop:
		for ;; callCount++ {
			select {
			case <-gen.stopSign:
				gen.handleStopSign(callCount)
				break Loop
			default:
			}
			gen.asyncCall()
			if gen.lps > 0 {
				select {
				case <-gen.cancelSign:
						gen.handleStopSign(callCount)
						break Loop
				case <-throttle:
				}
			}
		}
}

func (gen *myGenerator) handleStopSign(callCount uint64)  {
	// 先把信号值置成1
	gen.cancelSign = 1
	log.Println("Closing result channel...")
	// 还要记得关闭resultChan
	close(gen.resultCh)
	gen.endSign <- callCount
}

// 这个调用过程就是这个载荷发生器的核心之一了
// 1.首先要得从gorouite池里拿go票，结尾还要还回来
// 2.其次，要考虑调用过程中可能会发生些不可预料的错误造成程序恐慌，导致程序运行不下去
// 3.这里面省略掉的细节是为了不另起一个goroutie实现这个调用
// 3.1 本来是用两个select做的，其中一个case异步接收interact方法的结果，另一个case同时做超时处理
// 3.2 第一个改进是把interact收回到第一个case里面 <-func <-chan *lib.rawResponse {} 丢到匿名函数里面执行
// 3.3 但case的特点是会导致这个方法执行太久（已经执行了），再区考虑第2个case,这样时间就不准了
// 3.4 所以用到了time.AfterFunc函数（之前是time.After函数）
// 3.5 既然用了AfterFunc函数，这是超时后的判断，但是如果我们在不超时情况下跑了interact方法，跑着跑着，时间到了，也会执行这个超时
// 3.6 所以，不超时跑方法的时候，就需要先Stop掉这个超时（持续器）

func (gen *myGenerator) asyncCall() {
	gen.tickets.Take()
	go func() {
		// 调用过程中的未预料到的错误or调用器自身的错误
		// 不能让运行时的恐慌外泄并影响到正常流程的进行
		defer func() {
			if p := recover(); p !=nil {
				err, ok := interface{}(p).(error)
				var buff bytes.Buffer
				buff.WriteString("Async Call Panic!(")
				if ok {
					buff.WriteString("eror: ")
					buff.WriteString(err.Error())
				} else {
					buff.WriteString("clue: ")
					buff.WriteString(fmt.Sprintf("%v", p))
				}
				buff.WriteString(")")
				errMsg := buff.String()
				log.Fatalln(errMsg)
				result := &lib.CallResult{
					Id: -1,
					Code: lib.CODE_FATAL_CALL,
					Msg: errMsg}
				gen.SendResult(result)
				}
			}()

		rawRequest := gen.caller.BuildRequest()
		var timeout bool
		timer := time.AfterFunc(gen.timeoutNs, func() {
			timeout = true
			result := &lib.CallResult{
				Id: rawRequest.Id,
				Request: rawRequest,
				Code: lib.CODE_WARNING_CALL_TIMEOUT,
				Msg: fmt.Sprintf("Timeout! (expected: < %v)", gen.timeoutNs)}
			gen.senResult(result)
		})
		rawResponse := gen.interact(&rawRequest)
		if !timeout {
			timer.Stop()
			var result *lib.CallResult
			if rawResponse.Error != nil {
				result = &lib.CallResult{
					Id: rawRequest.Id,
					Request:rawRequest,
					Code:lib.CODE_ERROR_CALL,
					Msg:rawResponse.Error.Error(),
					Elapse:rawResponse.Elapse}
			} else {
				result = gen.caller.CheckResponse(rawResponse, *rawResponse)
				result.Elapse = rawResponse.Elapse
			}
			gen.sendResult(result)
		}
		gen.tickets.Return()
	}()
}

// 这个start也做了很多东西
// 1.设定throttle
// 2.AfterFunc 虽然看起来是初始化停止信号，但实际做的是持续一段时间就over
// 3.初始化endSign通道、状态、调用计数
func (gen *myGenerator) Start() {
	log.Println("Starting load generator...")
	// 设定节流阀
	var throttle <-chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		log.Printf("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	// 初始化停止信号
	go func() {
		time.AfterFunc(gen.durationNs, func() {
			log.Println("Stopping load generator...")
			gen.stopSign <- 0
		})
	}()

	// 初始化完结信号通道
	gen.endSign = make(chan uint64, 1)
	// 初始化调用执行计数
	gen.callCount = 0
	// 设置已启动状态
	go func() {
		// 生成载荷
		log.Println("Generating loads...")
		gen.genLoad(throttle)

		// 接收调用执行计数
		callCount := <-gen.endSign
		gen.status = lib.STATUS_STOPPED
		log.Printf("Stopped. (callCount=%d)\n", callCount)
	}()
}

// stop方法要考虑的就是前置检查了
func (gen *myGenerator) Stop() (uint64, bool) {
	if gen.stopSign == nil {
		return 0, false
	}
	if gen.status != lib.STATUS_STARTED {
		return 0, false
	}
	gen.status = lib.STATUS_STOPPED
	gen.stopSign <- 1
	callCount := <-gen.endSign

	return callCount, true
}

func (gen *myGenerator) interact(rawReq *lib.RawRequest) *lib.RawResponse {
	if rawReq == nil {
		return &lib.RawResponse{Id: -1, Error: errors.New("Invalid raw request.")}
	}
	start := time.Now().Nanosecond()
	resp, err := gen.caller.Call(rawReq.Request, gen.timeoutNs)
	end := time.Now().Nanosecond()
	elapsedTime := time.Duration(end - start)
	var rawResp lib.RawResponse
	if err != nil {
		errMsg := fmt.Sprintf("Sync Call Error: %s.", err)
		rawResp = lib.RawResponse{
			Id:     rawReq.Id,
			Error:    errors.New(errMsg),
			Elapse: elapsedTime}
	} else {
		rawResp = lib.RawResponse{
			Id:     rawReq.Id,
			Response:   resp,
			Elapse: elapsedTime}
	}
	return &rawResp
}
