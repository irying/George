# Go的练手项目
手写见过的各种项目，慕课网、《Go并发编程实战》、博客等项目

#### 1.loadgen -- 载荷发生器，测试TPS跟QPS  || abtest

《实战》书里面这个loadgen有点复杂，还不如模仿ab压测工具写个。

1）有多少个并发就有多少个goroutine

2) 给个time.AfterFunc函数，请求超时当作失败

3)给个p:=recover()函数，恢复进程

4)再给个通道，byte字节，表示停止信号，终止压测

#### 2.oneway -- 缓冲通道跟非缓冲通道的区别，channel使用  

#### 3.pipeline -- 外部排序（归并排序）  

#### 4.set -- 实现Go语言没有的set  

#### 5.mutux_cond -- 文件读写  

#### 6.pubsub -- 消息订阅和推送

#### 7.http代理

可以说是最简单的go工具了，关键在于这个转发函数。一个io.copy，copy完就去关闭对端的写端，为什么关闭对端的写端，就表示对端不再接收你的发送了。参考知乎问答https://www.zhihu.com/question/48871684/answer/113135138[胡宇光的回答]

> 从tcp协议本身来说，你是无法知道对方到底是调用了close()还是调用了shutdown(send)的，os的tcp协议栈也不知道。因此此时是否要close取决于你的应用。通常来说如果对方调用的是close，那么你也可以close。否则你不能close，例如对方发送一个包给你并shutdown write然后调用recv，这时候你还可以返回一个或多个包，连接此时处于半关闭状态，可以一直持续。这么做的客户端不多（connect, send, shtudown(send), recv();），但的确有，而且是完全合法的。
>
> 对于proxy来说，正确的做法是透传双方的行为。**因此，当你read(client_side_socket)返回0时，你应该对另外一端调用shutdown(server_side_socket, send)，这样服务器就会read返回0，你透明的传递了这个行为。**那么作为proxy，**你什么时候才能close呢？client_socket和server_socket上read都返回了0，或者有任何一方返回了-1时你可以close。**当然你也可以考虑设置一个超时时间，如果线路上超过5分钟没有数据你就断开，但这是另一个维度的问题。

todo:参考// https://golang.org/src/net/http/server.go#L2274 完善下断线重连

```go
func forward(destination net.Conn, source net.Conn) {
	io.Copy(destination, source)
	log.Printf("relay: done copying from %v to %v\n", 		     source.LocalAddr().String(), destination.RemoteAddr().String())
	tcpConnection := destination.(*net.TCPConn)
	tcpConnection.CloseWrite()
}
```

#### 8.context使用

1）如果是多个阶段协同完成一件事，则用到sync.waitGroup包，这个情况是等所有阶段结束，事情才算做完。

2）但如果是不相干的事情，我们想随时停止某个进程，比如为了实现一种需求：下游消费上游的数据，**当下游不需要数据时，上游能够停止生产。**这时就需要用到context，用cancel函数随意取消某个阶段。

