package mutexCond

import (
	"os"
	"sync"
)

type Data []byte

type DataFile interface {
	Read() (rsn int64, data Data, error error)
	Write(data Data) (wsn int64, error error)
	Rsn() int64
	Wsn() int64
	DataLen() uint32
}

type myDataFile struct {
	file *os.File
	fmutex sync.RWMutex //被用于文件的读写锁
	rcond *sync.Cond // 读操作需要的条件变量
	wmutex sync.Mutex // 写操作需要用到的互斥锁
	rmutex sync.Mutex // 读操作的互斥锁
	woffset int64
	roffset int64
	dataLen uint32
}
