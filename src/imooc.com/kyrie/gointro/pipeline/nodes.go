package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
)

func ArraySource( a ...int ) <- chan int {
	out := make(chan int)
	go func() {
		for _, v := range a{
			out<- v
		}
		close(out)
	}()
	return  out
}

func InMemorySort(in <-chan int) <-chan int  {
	out := make(chan int)
	go func() {
		a := [] int{}
		// save in memory
		for v := range in {
			a = append(a, v)
		}
		//sort
		sort.Ints(a)

		// output
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return  out
}

func ReaderSource(reader io.Reader) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err !=nil {
				break
			}
		}
		close(out)
	}()
	
	return  out
}
