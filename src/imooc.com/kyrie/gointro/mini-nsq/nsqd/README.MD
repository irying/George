nsqd的基本属性（目前）

```Go
type NSQD struct {
	clientIDSequence int64
	sync.RWMutex
	opts atomic.Value
	startTime time.Time
	httpListener net.Listener
	tcpListener net.Listener
	waitGroup waitwraper.WaitGroupWrapper

	topicMap map[string]*Topic
}
```

其中配置选项集中在options中

```go
type Options struct {
	ID int64 `flag:"node-id" cfg:"id"`
	TCPAddress string
	HTTPAddress string
	LookupdTCPAddress []string
}
```

### 第2个任务：起个http服务监听http请求

```go
httpListener, err := net.Listen("tcp", n.getOptions().HTTPAddress);
	if err != nil {
		// TODO 日志
	}
	ctx := &context{n}
	httpServer := newHTTPServer(ctx)
	n.waitGroup.Wrap(func() {
		http_api.Serve(httpListener,httpServer)
	})
```

主要看看httpServer结构，引入个httprouter包处理所有相关http路由(具体使用参考http://oopsguy.com/2017/10/20/build-api-server-with-httprouter-in-golang/)

```go
type httpServer struct {
	ctx *context
	router http.Handler
}

func newHTTPServer(ctx *context) *httpServer {
	router := httprouter.New()
	s := &httpServer{
		ctx:ctx,
		router:router,
	}
	...

	return s
}
```



### 第3个任务很有意思：queueScanLoop() 扫描和处理InFightQueue和DeferredQueue。

配置项增加几个跟此任务相关的配置

主要目的是扫描，但在扫描之前做了很多琐碎的工作，像启动多少个goroutine(n.poolSize)去扫描，调整goroutine数量；像处理中断的情况，比如定时器时间到了或者nsqd挂了退出；还有结果放到一个通道(responseCh)里面返回，根据redis的随机算法处理。

这些工作没做好，是不够完整和谨慎的。

#### 1.分配grouting，相当于worker

```go
func (n *NSQD) resizePool(num int, workCh chan *Channel, responseCh chan bool, closeCh chan int)  {
   // todo it confused me
   idealPoolSize := int(float64(num) * 0.25)
   if idealPoolSize < 1{
      idealPoolSize = 1
   } else if idealPoolSize > n.getOptions().QueueScanWorkerPoolMax{
      idealPoolSize = n.getOptions().QueueScanWorkerPoolMax
   }

   for  {
      if idealPoolSize == n.poolSize {
         break
      } else if idealPoolSize < n.poolSize{
         closeCh <- 1
         n.poolSize--
      } else {
         n.waitGroup.Wrap(func() {
            n.queueScanWorker(workCh, responseCh, closeCh)
         })
         n.poolSize ++
      }
   }

}
```



#### 2.处理中断的情况

```
select {
		case <-workTicker.C:
			if len(channels) == 0 {
				continue
			}
		case <-refreshTicker.C:
			channels = n.channels()
			n.resizePool(len(channels), workCh, responseCh, closeCh)
			continue
		case <-n.exitChan:
			goto exit
		}
```



#### 3.类似于redis删除过期键那种随机算法，随机扫描channel

```go
loop:
		for _, i := range util.UniqRands(num, len(channels)){
			workCh<-channels[i]
		}
		
		numDirty := 0
		for i:=0; i < num; i++ {
			if <-responseCh {
				numDirty++
			}
		}

		if float64(numDirty)/float64(num) > n.getOptions().QueueScanDirtyPercent {
			goto loop
		}
```

1. `workTicker`定时器触发扫描流程。 nsqd采用了Redis的probabilistic expiration算法来进行扫描。首先从所有Channel中随机选取部分Channel，然后遍历被选取的Channel，投到`workerChan`中，并且等待反馈结果， 结果有两种，dirty和非dirty，如果dirty的比例超过配置中设定的`QueueScanDirtyPercent`，**那么不进入休眠， 继续扫描，如果比例较低，则重新等待定时器触发下一轮扫描。**这种机制可以在保证处理延时较低的情况下减少对CPU资源的浪费。
2. `refreshTicker`**定时器触发更新Channel列表流程**。 这个流程比较简单，先获取一次Channel列表， 再调用`resizePool`重新分配worker。

 获取channel出来

```Go
func (n *NSQD) channels() []*Channel {
	var channels []*Channel
	n.RLock()
	for _, t := range n.topicMap{
		t.RLock()
		for _, c := range t.channelMap{
			channels = append(channels, c)
		}
		t.RUnlock()
	}
	n.RUnlock()
	
	return channels
}
```

#### 琐碎的工作有了，到正事这边

每个worker做的事

`queueScanWorker`接收`workCh`发来的消息，处理，并且通过`responseCh`反馈消息。收到`closeCh`时则关闭。 **由于所有worker都监听相同的`closeCh`，所以当向`closeCh`发送消息时，随机关闭一个worker**。 且由于`workCh`和`closeCh`的监听是串行的，所以不存在任务处理到一半时被关闭的可能。 这也是nsq中优雅关闭gochannel的的一个例子。

```Go
func (n *NSQD) queueScanWorker(workCh chan *Channel, responseCh chan bool, closeCh chan int)  {
	for {
		select {
		case c:= <-workCh:
			now := time.Now().UnixNano()
			dirty := false
			if c.processInFlightQueue(now) {
				dirty = true
			}
			if c.processDeferredQueue(now) {
				dirty = true
			}
			responseCh <- dirty
		case <-closeCh:
			return
		}
	}
}
```





消息分发可靠性

会把发送到客户端的信息存到inFlightMessages和inFlightPQ，当一条消息被发送给client之后，如果在指定时间内仍没有收到客户端的响应的话，NSQ会把本次发送当做失败处理，消息还在队列中。如果收到客户端响应，FIN信号，消息就完成，从队列跟map中删除



nsqd启动的时候有个计划任务，定期扫描这些队列，进行发送。如果有多个消费者的话