package handle

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"trex_online/src/match"
)

type OnKeyReq struct {
	UserId  string `json:"userId"`
	KeyCode int32  `json:"keyCode"`
}

// todo 区分玩家
func OnKeyUp(s socketio.Socket, reqString string) {
	req := new(OnKeyReq)
	if err := json.Unmarshal([]byte(reqString), req); err != nil {
		return
	}
	room := match.MEngine.GetUserRoom(req.UserId)
	if room == nil {

	} else {
		room.Broadcast(reqString, "keyUp")
	}
	return
}

func OnKeyDown(s socketio.Socket, reqString string) {
	req := new(OnKeyReq)
	if err := json.Unmarshal([]byte(reqString), req); err != nil {
		return
	}
	room := match.MEngine.GetUserRoom(req.UserId)
	if room == nil {

	} else {
		room.Broadcast(reqString, "keyDown")
	}
	return
}
