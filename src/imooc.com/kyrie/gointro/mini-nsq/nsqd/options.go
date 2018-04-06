package nsqd

import (
	"os"
	"log"
	"crypto/md5"
	"io"
	"hash/crc32"
	"time"
)

type Options struct {
	ID int64 `flag:"node-id" cfg:"id"`
	TCPAddress string
	HTTPAddress string
	LookupdTCPAddress []string

	// 队列扫描
	QueueScanInterval time.Duration
	QueueScanRefreshInterval time.Duration
	QueueScanSelectionCount int
	QueueScanWorkerPoolMax int
	QueueScanDirtyPercent float64
}

func NewOptions() *Options  {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	// TODO 这个id计算很有意思
	h := md5.New()
	io.WriteString(h, hostname)
	defaultID := int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)

	return  &Options{
		ID:defaultID,
		TCPAddress:"0.0.0.0:4150",
		HTTPAddress:"0.0.0.0:4151",
		LookupdTCPAddress:make([]string, 0),
	}
}
