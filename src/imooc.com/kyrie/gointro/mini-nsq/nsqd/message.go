package nsqd

const (
	MsgIDLength = 16
	minValidMsgLength = MsgIDLength + 8 + 2
)

type MessageID [MsgIDLength]byte

type Message struct {
	ID MessageID
	Body []byte
	Timestamp int64
	Attempts uint16
	index int


	clientID int64
}