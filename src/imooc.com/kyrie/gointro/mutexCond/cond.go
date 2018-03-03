package mutexCond

type Data []byte

type DataFile interface {
	Read() (rsn int64, data Data, error error)
	Write(data Data) (wsn int64, error error)
	Rsn() int64
	Wsn() int64
	DataLen() uint32
}


