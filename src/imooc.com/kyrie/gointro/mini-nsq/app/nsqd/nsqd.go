package nsqd

import (
	"imooc.com/kyrie/gointro/mini-nsq/nsqd"
	"github.com/judwhite/go-svc/svc"
	"syscall"
	"log"
)

type program struct {
	nsqd *nsqd.NSQD
}

func main()  {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Start() error {
	opts := nsqd.NewOptions()
	nsq := nsqd.New(opts);
	nsq.Main()
	p.nsqd = nsq
	return nil
}