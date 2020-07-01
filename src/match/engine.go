package match

import (
	"log"
	"sync"
)

const (
	RoomCap       = 2
	MaxMatchQueue = 1000
)

var MEngine *Engine

type Engine struct {
	roomMap     sync.Map
	userRoomMap sync.Map
	matchQueue  chan *User
	latestRoom  *Room
	matchLock   chan struct{}
	roomCap     int32
}

func NewMatchEngine() *Engine {
	engine := &Engine{
		roomMap:    sync.Map{},
		matchQueue: make(chan *User, MaxMatchQueue),
		matchLock:  make(chan struct{}, 1),
		roomCap:    RoomCap,
	}
	engine.NewRoom()
	return engine
}

func (e *Engine) NewRoom() {
	room := NewRoom()
	e.latestRoom = room
	e.roomMap.Store(room.Id, room)
}

func (e *Engine) JoinRoom(user *User) {
	e.matchQueue <- user
}

func (e *Engine) GetUserRoom(userId string) *Room {
	if roomIn, ok := e.userRoomMap.Load(userId); !ok {
		return nil
	} else {
		return roomIn.(*Room)
	}
}

func (e *Engine) LeaveRoom(userId string) {
	if roomIn, ok := e.userRoomMap.Load(userId); !ok {
		return
	} else {
		room := roomIn.(*Room)
		if !room.Full {
			user := room.GetUser(userId)
			if user != nil {
				e.matchLock <- struct{}{}
				room.Len -= 1
				user.LeaveRoom(room.Name)
				room.removeUser(userId)
				e.userRoomMap.Delete(userId)
				<-e.matchLock
				log.Println("leave")
			}
		}
		return
	}
}

func (e *Engine) MatchLoop() {
	for {
		select {
		case user := <-e.matchQueue:
			// 如果当前房间已满，则创建新的房间
			if e.latestRoom.Full {
				e.NewRoom()
			} else if user := e.latestRoom.GetUser(user.Id); user != nil {
				log.Println("already")
				break
			}
			e.matchLock <- struct{}{}
			// 用户进入房间
			user.GetInRoom(e.latestRoom.Id, e.latestRoom.Name)
			e.latestRoom.AddUser(user)
			e.userRoomMap.Store(user.Id, e.latestRoom)
			e.latestRoom.Len += 1
			// 是否已满
			log.Println(e.latestRoom.Len)

			if e.latestRoom.Len == e.latestRoom.Cap {
				e.latestRoom.Full = true
				// todo 开始对战
				go e.latestRoom.RunBattle()
			}
			<-e.matchLock
		}
	}
}

func (e *Engine) Run() {
	go e.MatchLoop()
}
