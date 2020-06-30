package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"net/http"
	"trex_online/src/match"
)

func main() {

	// 启动 match engine
	matchEngine := match.NewMatchEngine()
	matchEngine.Run()
	for i := 0; i < 10000; i++ {
		matchEngine.JoinRoom(&match.User{})
	}

	var (
		err          error
		SocketServer *socketio.Server
	)
	if SocketServer, err = socketio.NewServer(nil); err != nil {
		panic("ws server error" + err.Error())
	}

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
