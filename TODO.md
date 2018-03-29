TODO

一、nsqlookup

二、nsqd

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

- [ ] ​

      ​