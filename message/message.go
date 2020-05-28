package message

const (
	Plain int = iota
	Heartbeat
)

type Message struct {
	MessageType int
	Content     string
	Owner       string
}

func (m *Message) IsHeartBeat() bool {
	return m.MessageType == Heartbeat
}
