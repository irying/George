package main

import (
	"github.com/nsqio/go-nsq"
	"time"
	"fmt"
	"utils/waitwraper"
)

func main() {
	var wg waitwraper.WaitGroupWrapper
	//接受消息
	consume()
	//分别向不同的服务节点发送消息
	wg.Wrap(func(){ produce("node1","localhost:4150")})
	wg.Wrap(func (){produce("node2","localhost:4152")})

	wg.Wait()
}
func produce(tag string,addr string) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second*5)
		p.Publish("test", []byte(tag+":"+time.Now().String()))
	}
}
func consume() {
	config := nsq.NewConfig()
	//注意MaxInFlight的设置，默认只能接受一个节点
	config.MaxInFlight=2
	c, err := nsq.NewConsumer("test", "consum", config)
	if err != nil {
		panic(err)
	}
	hand := func(msg *nsq.Message) error{
		fmt.Println(string(msg.Body))
		return nil
	}
	c.AddHandler(nsq.HandlerFunc(hand))
	if err:= c.ConnectToNSQLookupd("localhost:4161");err!=nil{
		fmt.Println(err)
	}
}
