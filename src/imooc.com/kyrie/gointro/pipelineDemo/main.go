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
	p1 := pipeline.InMemorySort(pipeline.ArraySource(3, 2, 6, 7, 4))
	p2 := pipeline.InMemorySort(pipeline.ArraySource(7, 12, 8, 9, 0, 5, 8, 10))
	p := pipeline.Merge(p1, p2)

	for v := range p {
		fmt.Println(v)
	}
}
