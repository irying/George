package main

import (
	"bufio"
	"os"
	"github.com/nsqio/go-nsq"
	"fmt"
)

var producer *nsq.Producer

func main() {
	ip1 := "127.0.0.1:4150"
	ip2 := "127.0.0.1:4152"
	InitProducer(ip1)
	running := true

	// 读取控制台输入
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}

		for err := Publish("test", command); err != nil; err = Publish("test", command) {
			ip1, ip2 = ip2, ip1
			InitProducer(ip1)
		}
		// 关闭
		producer.Stop()
	}


}
func InitProducer(ip string) {
	var err error
	producer, err = nsq.NewProducer(ip, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
}
// Producer不能发布(Publish)空message，否则会导致panic
func Publish(topic string, message string) error {
	var err error
	if producer != nil {
		if message == ""{
			return nil
		}

		err = producer.Publish(topic, []byte(message))
		return err
	}
	return fmt.Errorf("producer is nil", err)
}
