package msg

import socketIo "github.com/googollee/go-socket.io"

const (
	MaxMsgQueue = 1000
)

var Message *BroadCast

type BroadCast struct {
	sServer  *socketIo.Server
	msgQueue chan *Msg
}

type Msg struct {
	Msg   string
	Room  string
	Event string
}

func NewBroadCast(server *socketIo.Server) *BroadCast {
	return &BroadCast{
		sServer:  server,
		msgQueue: make(chan *Msg, MaxMsgQueue),
	}
}

func (b *BroadCast) AddMsg(msg *Msg) {
	b.msgQueue <- msg
}

func (b *BroadCast) Run() {
	go b.MsgLoop()
}

func (b *BroadCast) MsgLoop() {
	for {
		select {
		case msg := <-b.msgQueue:
			b.sServer.BroadcastTo(msg.Room, msg.Event, msg.Msg)
		}
	}
}
