package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"net/http"
	"trex_online/src/handle"
	"trex_online/src/match"
	"trex_online/src/msg"
)

func main() {

	var (
		err          error
		SocketServer *socketio.Server
	)
	if SocketServer, err = socketio.NewServer(nil); err != nil {
		panic("ws server error" + err.Error())
	}

	msg.Message = msg.NewBroadCast(SocketServer)
	msg.Message.Run()
	// 启动 match engine
	match.MEngine = match.NewMatchEngine()
	match.MEngine.Run()

	SocketServer.On("match", handle.Match)
	SocketServer.On("getInRoom", handle.GetInRoom)
	SocketServer.On("leaveRoom", handle.LeaveRoom)
	SocketServer.On("keyDown", handle.OnKeyDown)
	SocketServer.On("keyUp", handle.OnKeyUp)

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		SocketServer.ServeHTTP(w, r)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%s", "0.0.0.0", "8080"), nil)
}
