package match

import socketIo "github.com/googollee/go-socket.io"

type User struct {
	Id string
	s  socketIo.Socket
}

func (u *User) GetInRoom(roomId int32, roomName string) {
	u.s.Join(roomName)
	u.s.Emit("getInRoom", roomId)
}

func (u *User) LeaveRoom(roomName string) {
	u.s.Leave(roomName)
}

func (u *User) SetSocket(socket socketIo.Socket) {
	u.s = socket
}
