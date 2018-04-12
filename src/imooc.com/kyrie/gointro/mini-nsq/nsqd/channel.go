package nsqd

import (
	"sync"
	"sync/atomic"
	"imooc.com/kyrie/gointro/mini-nsq/internal/pqueue"
	"time"
	"errors"
)

type Channel struct {
	sync.RWMutex
	topicName string   // 所属topic的名字
	name      string   // channel的名字
	ctx       *context // 上下文对象（NSQD）

	clients        map[int64]Consumer

	messageMsgChan chan *Message
	exitFlag int32
	exitMutex sync.RWMutex

	// 队列部分来了
	deferredMessages map[MessageID]*pqueue.Item
	deferredPQ pqueue.PriorityQueue
	deferredMutex sync.Mutex

	inFlightMessages map[MessageID]*Message
	inFlightPQ inFlightPqueue
	inFlightMutex sync.Mutex
}

func (c *Channel) Exiting() bool {
	return atomic.LoadInt32(&c.exitFlag) == 1
}

func (c *Channel) processInFlightQueue(t int64)  {
	c.exitMutex.RLock()
	defer c.exitMutex.RUnlock()

	if c.Exiting() {
		return false
	}

	dirty := false
	for {
		c.inFlightMutex.Lock()
		msg, _:= c.inFlightPQ.PeekAndShift(t)
		c.inFlightMutex.Unlock()

		if msg == nil {
			goto exit
		}
		dirty = true
		_, err := c.popInFlightMessage(msg.clientID, msg.ID)
		if err != nil {
			goto exit
		}
		// atomic.AddInt64(&c.timeoutCount, 1)
		c.RLock()
		// todo customer部分
		// client, ok := c.clients[msg.clientID]
		c.RUnlock()
		//if ok {
		// client.TimedOutMessage()
		//}
		c.put(msg)
	}
	exit:
	return dirty
}

// popInFlightMessage atomically removes a message from the in-flight dictionary
func (c *Channel) popInFlightMessage(clientID int64, id MessageID) (*Message, error) {
	c.inFlightMutex.Lock()
	msg, ok := c.inFlightMessages[id]
	if !ok {
		c.inFlightMutex.Unlock()
		return nil, errors.New("ID not in flight")
	}
	if msg.clientID != clientID {
		c.inFlightMutex.Unlock()
		return nil, errors.New("client does not own message")
	}
	delete(c.inFlightMessages, id)
	c.inFlightMutex.Unlock()
	return msg, nil
}