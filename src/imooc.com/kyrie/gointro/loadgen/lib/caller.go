package lib

import "time"

type Caller interface {
	BuildRequest() RawRequest
	// 调用
	Call(request []byte, timeoutNs time.Duration) ([]byte, error)
	CheckResponse(rawRequest RawRequest, rawResponse RawResponse) *CallResult
}