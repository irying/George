package pipeline

import "sort"

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
