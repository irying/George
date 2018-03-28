package main

import (
	"time"
	"github.com/nsqio/go-nsq"
	"fmt"
)

type ConsumerT struct {}

func main() {
	InitConsumer("test", "test-channel", "127.0.0.1:4161")
	for  {
		time.Sleep(time.Second * 10)
	}
}

func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	return nil
}

func InitConsumer(topic string, channel string, address string) {
	config := nsq.NewConfig()
	config.LookupdPollInterval = time.Second
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		panic(err)
	}

	consumer.SetLogger(nil, 0)
	consumer.AddHandler(&ConsumerT{})

	if err := consumer.ConnectToNSQLookupd(address); err !=nil {
		panic(err)
	}
}
