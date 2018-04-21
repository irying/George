package main

import (
	"context"
	"fmt"
	"time"
)

var key string="name"

func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	value1Ctx := context.WithValue(ctx, key, "monitor first")
	value2Ctx := context.WithValue(ctx, key, "monitor second")
	value3Ctx := context.WithValue(ctx, key, "monitor thrid")
	go watch(value1Ctx)
	go watch(value2Ctx)
	go watch(value3Ctx)
	time.Sleep(10 * time.Second)
	fmt.Println("ok, let' s stop")
	cancel()
	time.Sleep(5 * time.Second)


}
func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key),"stop")
			return
		default:
			fmt.Println(ctx.Value(key),"is running")
			time.Sleep(2 * time.Second)
		}
	}
}
