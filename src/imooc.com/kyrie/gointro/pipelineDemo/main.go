package main

import (
	"imooc.com/kyrie/gointro/pipeline"
	"fmt"
)

func main() {
	//p := pipeline.ArraySource(3, 2, 6, 7, 4)
	//for {
	//	if num, ok := <-p; ok {
	//		fmt.Println(num)
	//	} else {
	//		break
	//	}
	//}
	p := pipeline.InMemorySort(pipeline.ArraySource(3, 2, 6, 7, 4))
	for v := range p {
		fmt.Println(v)
	}
}
