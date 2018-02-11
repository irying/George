package main

import (
	"imooc.com/kyrie/gointro/pipeline"
	"fmt"
	"os"
)

func main() {
	const filename = "small.in"
	const count = 50
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p := pipeline.RandomSource(count)
	pipeline.WriterSink(file, p)

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(file)
	for v := range p{
		fmt.Println(v)
	}
}

func MergeDemo() {//p := pipeline.ArraySource(3, 2, 6, 7, 4)
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