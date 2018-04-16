package main

import (
	"flag"
	"time"
	"sync"
	"net/http"
	"fmt"
)

type Work struct {
	Requests int
	Concurrency int
	Timeout int
	Url string
	results chan *Result
	start time.Time
	end time.Time
}

type Result struct {
	// request *Request
	// response *Response
	Duration time.Duration
}

func main()  {
	var (
		requests int
		concurrency int
		timeout int
		url string
	)

	flag.IntVar(&requests, "n", 1000, "")
	flag.IntVar(&concurrency, "c", 100, "")
	flag.IntVar(&timeout, "s", 10, "")
	flag.StringVar(&url, "url", "http://www.baidu.com", "")
	flag.Parse()
	w := Work{
		Requests:requests,
		Concurrency:concurrency,
		Timeout:timeout,
		Url:url,
	}
	w.run()
}

func (w *Work)run()  {
	w.results = make(chan *Result, w.Requests)
	w.start = time.Now()
	w.runWorkers();
	w.end = time.Now()
	w.printResult()
}

func (w *Work)runWorkers() {
	var wg sync.WaitGroup
	wg.Add(w.Concurrency)
	for i := 0; i <= w.Concurrency; i++  {
		go func() {
			defer wg.Done()
			w.runWorker()
		}()
	}
	wg.Wait()
	close(w.results)
}

func (w *Work)runWorker()  {
	num := w.Requests/w.Concurrency
	client := &http.Client{Timeout:time.Duration(w.Timeout) * time.Second}
	for i := 0; i <= num ; i++ {
		w.sendRequest(client)
	}
}

func (w *Work)sendRequest(client *http.Client) {
	request, err := http.NewRequest("GET", w.Url, nil)
	if err!=nil {
		// todo log
	}
	start := time.Now()
	client.Do(request)
	end := time.Now()
	w.results <- &Result{
		Duration:end.Sub(start),
	}
}

func (w *Work) printResult() {
	sum := 0.0
	num := float64(len(w.results))

	for result := range w.results {
		sum += result.Duration.Seconds()
	}

	rps := int(num / w.end.Sub(w.start).Seconds())
	tpr := sum/ num * 1000

	fmt.Printf("Requests per second:\t%d [#/sec]\n", rps)
	fmt.Printf("Time per request:\t%.3f [ms]\n", tpr)

}

