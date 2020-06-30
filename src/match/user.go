package match

import socketIo "github.com/googollee/go-socket.io"

type User struct {
	Id string
	s  socketIo.Socket
}
