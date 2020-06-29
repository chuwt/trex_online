package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
)

func main() {
	var err error
	if SocketServer, err = socketio.NewServer(nil); err != nil {
		panic("ws server error" + err.Error())
	}

	SocketServer.On("connect_", Connect_)
	SocketServer.On("match", Match)
	SocketServer.On("room", GetRoom)
	SocketServer.On("leave", LeaveRoom)
	SocketServer.On("keyup", KeyUp)
	SocketServer.On("keydown", KeyDown)

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
		SocketServer.ServeHTTP(w, r)
	})

	go runMatch()

	http.ListenAndServe(fmt.Sprintf("%s:%s", "0.0.0.0", "8080"), nil)
}

var SocketServer *socketio.Server

func Connect_(s socketio.Socket, req string) {
	s.Join("game")
	return
}

var (
	Start              = false
	ReadyCount         = 0
	ReadyMap           = make(map[string]bool)
	Speed      float64 = 0
	MaxSpeed float64 = 13
	SpeedStep float64= 0.01

	MatchChan = make(chan User, 200)

	// userId: socket
	MatchedMap = sync.Map{}
)

func Match(s socketio.Socket, req string) {
	if _, ok := MatchedMap.Load(req); ok {
		return
	}
	MatchChan <- User{
		s:    s,
		Name: req,
	}
}

func LeaveRoom(s socketio.Socket, req string) {
	waitRoomIn, ok := MatchedMap.Load(req)
	if !ok {
		return
	}
	waitRoom := waitRoomIn.(*Room)
	if waitRoom.Finished {
		// 对局结束才能离开房间
		MatchedMap.Delete(req)
	}
}

func GetRoom(s socketio.Socket, req string) {
	if in, ok := MatchedMap.Load(req); !ok {
		return
	} else {
		room := in.(*Room)
		s.Join(room.Name)
		s.Emit("roomRes", room.Name)
	}
}

type User struct {
	s    socketio.Socket
	Name string
}

type Room struct {
	Name   string
	isFull bool
	//People   int32 // 人数
	UserList []User
	FinishedChan chan bool
	Finished bool
}

func runMatch() {
	waitRoom := new(Room)
	for {
		select {
		case user := <-MatchChan:
			if waitRoom.isFull {
				waitRoom = new(Room)
			} else {
				// 防止重入房间
				if _, ok := MatchedMap.Load(user.Name); ok {
					continue
				}
				MatchedMap.Store(user.Name, waitRoom)
				waitRoom.UserList = append(waitRoom.UserList, user)
				if len(waitRoom.UserList) == 2 {
					log.Println("匹配成功")
					waitRoom.Name = CreateRandomString(12)
					waitRoom.isFull = true
					// 房间满员，开始对战
					for _, user := range waitRoom.UserList {
						user.s.Join(waitRoom.Name)
						MatchedMap.Store(user.Name, waitRoom)
					}
					SocketServer.BroadcastTo(waitRoom.Name, "roomRes", waitRoom.Name)
					go func(waitRoom *Room) {
						Speed = 3
						for {
							select {
							case <-waitRoom.FinishedChan:
								// 退出房间
								close(waitRoom.FinishedChan)
								waitRoom.Finished = true
								log.Println(waitRoom.Name, "结束")
								return
							default:
								if Speed <= MaxSpeed {
									Speed += SpeedStep
								}
								SocketServer.BroadcastTo(waitRoom.Name, "speedRes", Speed)
							}
							time.Sleep(time.Second)
						}
						// todo 障碍生成
					}(waitRoom)
				} else {
					log.Println("等待匹配中")
				}
			}
		}
	}
}

func Ready(s socketio.Socket, req string) {
	fmt.Println(ReadyCount)
	s.BroadcastTo("game", "readyRes", req)
	if _, ok := ReadyMap[req]; !ok {
		ReadyCount++
		ReadyMap[req] = true
	}
	if ReadyCount == 2 && !Start {
		// 运行
		Start = true
		go func() {
			for {
				SocketServer.BroadcastTo("game", "speedRes", Speed)
				if Speed <= 13 {
					Speed += 0.01
				}
				time.Sleep(time.Millisecond * 200)
			}
		}()
	}
	return
}

func KeyUp(s socketio.Socket, req int) {
	s.BroadcastTo("game", "keyupRes", req)
	return
}

func KeyDown(s socketio.Socket, req int) {
	s.BroadcastTo("game", "keydownRes", req)
	return
}

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
