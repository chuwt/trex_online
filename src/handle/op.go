package handle

import (
	socketio "github.com/googollee/go-socket.io"
	"trex_online/src/match"
)

// todo 区分玩家
func OnKeyUp(s socketio.Socket, userId string) {
	room := match.MEngine.GetUserRoom(userId)
	if room == nil {

	} else {
		room.Broadcast("38", "keyUp")
	}
	return
}

func OnKeyDown(s socketio.Socket, userId string) {
	room := match.MEngine.GetUserRoom(userId)
	if room == nil {

	} else {
		room.Broadcast("38", "keyDown")
	}
	return
}
