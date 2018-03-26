package lib

import "time"

type ResultCode int
type GenStatus int

const (
	CODE_SUCCESS = 0
	CODE_WARNING_CALL_TIMEOUT ResultCode = 1001
	CODE_ERROR_CALL ResultCode = 1002
	CODE_ERROR_RESPONSE ResultCode = 1003
	CODE_ERROR_CALLEE ResultCode = 1004
	CODE_FATAL_CALL ResultCode = 1005
)

func GetResultCodePlain(code ResultCode) string {
	var codePlain string
	switch code {
	case CODE_SUCCESS:
		codePlain = "Success"
	case CODE_WARNING_CALL_TIMEOUT:
		codePlain = "Call Timeout Warning"
	case CODE_ERROR_CALL:
		codePlain = "Call Error"
	case CODE_ERROR_RESPONSE:
		codePlain = "Response Error"
	case CODE_ERROR_CALLEE:
		codePlain = "Callee Error"
	case CODE_FATAL_CALL:
		codePlain = "Call Fatal Error"
	default:
		codePlain = "Unknown result code"
	}
	return codePlain
}

const (
	STATUS_ORIGINAL GenStatus = 0
	STATUS_STARTED  GenStatus = 1
	STATUS_STOPPED  GenStatus = 2
)

type Generator interface {
	Start()
	// 第一个结果值代表已发载荷总数，且仅在第二个结果值为true时有效。
	// 第二个结果值代表是否成功将载荷发生器转变为已停止状态。
	Stop() (uint64, bool)
	Status() GenStatus
}

type RawRequest struct {
	Id      int64
	Request []byte
}

type RawResponse struct {
	Id int64
	Response []byte
	Error error
	Elapse time.Duration // 处理耗时
}

type CallResult struct {
	Id int64
	Request RawRequest
	Response RawResponse
	Code ResultCode
	Msg string
	Elapse time.Duration
}