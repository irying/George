lookupd的基本属性

```go
type Lookupd struct {
	sync.RWMutex
	opts *Options
	tcpListener net.Listener
	httpListener net.Listener
	waitGroup waitwraper.WaitGroupWrapper
}
```