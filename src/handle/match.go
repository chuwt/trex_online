package handle

import (
	socketio "github.com/googollee/go-socket.io"
	"trex_online/src/match"
)

func Match(s socketio.Socket, userId string) {
	user := &match.User{
		Id: userId,
	}
	user.SetSocket(s)
	match.MEngine.JoinRoom(user)
	return
}

func GetInRoom(s socketio.Socket, userId string) {
	room := match.MEngine.GetUserRoom(userId)
	if room == nil {
		// 不存在
		s.Emit("getInRoom", -1)
		return
	} else {
		user := room.GetUser(userId)
		if user != nil {
			user.SetSocket(s)
			// 订阅room
			user.GetInRoom(room.Id, room.Name)
			userList := room.GetRoomUsers()
			s.Emit("startBattle", userList)
		}
	}
	return
}

func LeaveRoom(s socketio.Socket, userId string) {
	match.MEngine.LeaveRoom(userId)
	return
}
