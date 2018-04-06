package nsqd

import (
	"sync"
	"sync/atomic"
	"time"
	"net"
	"utils/waitwraper"
	"imooc.com/kyrie/gointro/mini-nsq/internal/http_api"
	"imooc.com/kyrie/gointro/mini-nsq/internal/util"
)

type NSQD struct {
	clientIDSequence int64
	sync.RWMutex
	opts atomic.Value
	startTime time.Time
	httpListener net.Listener
	tcpListener net.Listener
	waitGroup waitwraper.WaitGroupWrapper
	poolSize int

	// todo 我打算用byte
	exitChan chan int

	topicMap map[string]*Topic
}

func New(opts *Options) *NSQD {
	n := &NSQD{
		startTime:time.Now(),
		topicMap:make(map[string]*Topic),
	}
	n.opts.Store(opts)
	return n
}

func (n *NSQD)getOptions() *Options  {
	return n.opts.Load().(*Options);
}

func (n *NSQD) Main()  {
	httpListener, err := net.Listen("tcp", n.getOptions().HTTPAddress);
	if err != nil {
		// TODO 日志
	}
	ctx := &context{n}
	httpServer := newHTTPServer(ctx)
	n.waitGroup.Wrap(func() {
		http_api.Serve(httpListener,httpServer)
	})
	n.waitGroup.Wrap(func() {
		n.queueScanLoop()
	})

}

func (n *NSQD) queueScanLoop()  {
	workCh := make(chan *Channel, n.getOptions().QueueScanSelectionCount)
	responseCh := make(chan bool, n.getOptions().QueueScanSelectionCount)
	closeCh := make(chan int)

	workTicker := time.NewTicker(n.getOptions().QueueScanInterval)
	refreshTicker := time.NewTicker(n.getOptions().QueueScanRefreshInterval)

	channels := n.channels()
	n.resizePool(len(channels), workCh, responseCh, closeCh)
	for  {
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

		num := n.getOptions().QueueScanSelectionCount
		if num > len(channels) {
			num = len(channels)
		}
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

	}

	exit:
	// n.logf("QUEUESCAN: closing")
	close(closeCh)
	workTicker.Stop()
	refreshTicker.Stop()

}

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