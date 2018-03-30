TODO

一、nsqlookup

当nsq集群中有多个nsqlookupd服务时，因为每**个nsqd都会向所有的nsqlookupd上报本地信息，因此nsqlookupd具有最终一致性**。

- [ ] tcp处理流程


- [ ] PING nsqd每隔一段时间都会向nsqlookupd发送心跳，表明自己还活着；
- [ ] IDENTITY 当nsqd第一次连接nsqlookupd时，发送IDENTITY，验证自己身份；
- [ ] REGISTER 当nsqd创建一个topic或者channel时，向nsqlookupd发送REGISTER请求，在nsqlookupd上更新当前nsqd的topic或者channel信息；
- [ ] UNREGISTER 当nsqd删除一个topic或者channel时，向nsqlookupd发送UNREGISTER请求，在nsqlookupd上更新当前nsqd的topic或者channel信息； 具体各个命令怎么执行，这里就不去分析了；需要提一点是，nsqd的信息是保存在registration_db这样的实例里面的；



二、nsqd 是一个守护进程，负责接收，排队，投递消息给客户端。

它可以独立运行，不过通常它是由 `nsqlookupd` 实例所在集群配置的（它在这能声明 topics 和 channels，以便大家能找到）。

- [ ] 1.github.com/judwhite/go-svc/svc

      ```
      svc框架启动一个service
      ```

- [ ] 2.配置文件，flag命令解析

- [ ] 3.json包解析磁盘中dat元数据

三、topic

四、channel

- [ ] channel在投递消息前， 会自增`msg.Attempts`，该变量用于保存投递尝试的次数。

- [ ] 在消息投递前会将`bufferedCount`置为1，在投递后置为0。该变量在`Depth`函数中被调用。

      ```
      func (c *Channel) Depth() int64 {
          return int64(len(c.memoryMsgChan)) + c.backend.Depth() + int64(atomic.LoadInt32(&c.bufferedCount))
      }
      ```

      `Deepth`函数返回内存，磁盘以及正在投递的消息数量之和，也就是尚未投递成功的消息数。

- [ ] messages are delivered at least once

- [ ] messages received are un-ordered

- [ ] ##### queueScanWorker，InFlightQueue 和 DeferredQueue

- [ ] NSQ使用priority queue和map来提高查找和使用消息的效率，但是可能会有潜在冗余和一致性问题。

- [ ] client可以通过设置timeout、sampleRate、maxAttemptCount等参数来过滤接收到的消息。