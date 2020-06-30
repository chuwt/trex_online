package msg

import socketIo "github.com/googollee/go-socket.io"

const (
	MaxMsgQueue = 1000
)

type BroadCast struct {
	sServer  *socketIo.Server
	msgQueue chan string
}

func NewBroadCast(server *socketIo.Server) *BroadCast {
	return &BroadCast{
		sServer:  server,
		msgQueue: make(chan string, MaxMsgQueue),
	}
}

func (b *BroadCast) AddMsg(msg string) {
	b.msgQueue <- msg
}

func (b *BroadCast) MsgLoop() {
	for {
		select {
		case msg := <-b.msgQueue:
			b.sServer.BroadcastTo("", "", msg)
		}
	}
}
