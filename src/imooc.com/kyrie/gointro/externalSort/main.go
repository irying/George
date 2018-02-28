package main

import (
	"imooc.com/kyrie/gointro/pipeline"
	"os"
	"bufio"
	"strconv"
)

func main() {
	p := CreateNetWorkPipeline("large.in", 800000000, 4)
	WriteToFile(p, "large.out")
	PrintFile("large.out")
}
// 1.source读文件
// 2.NetworkSink排好序写到chan
// 3.pipeline归并各个地址的
func CreateNetWorkPipeline(fileName string, fileSize int, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	pipeline.Init()
	
	sortAddr := []string{}
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		addr := ":" + strconv.Itoa(7000 + i)
		pipeline.NetworkSink(addr, source)
		sortAddr = append(sortAddr, addr)
	}

	sortResults := []<-chan int{}
	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}

	return pipeline.MergeN(sortResults)
}
