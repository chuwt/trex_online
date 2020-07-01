package match

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"trex_online/src/msg"
	"trex_online/src/util"
)

const (
	RoomId     int32 = 0
	MaxSpeed         = 13
	StartSpeed       = 3.0
)

var (
	SpeedStep = 0.01
)

type Room struct {
	Name string // 房间名称
	Id   int32

	Cap        int32    // 容量
	Len        int32    // 数量
	Finish     bool     // 对局是否完成
	Full       bool     // 是否满员
	people     sync.Map // 对局的用户
	peopleList []string // 用户列表

	CreateTime int64 // 时间
	FinishTime int64 // 完成时间

}

func NewRoom() *Room {
	roomId := util.Int32Incr(RoomId)
	room := &Room{
		Name:       util.CreateRandomString(8),
		Id:         roomId,
		Cap:        RoomCap,
		Len:        0,
		Finish:     false,
		people:     sync.Map{},
		CreateTime: time.Now().Unix(),
		FinishTime: 0,
	}
	return room
}

func (r *Room) GetUser(userId string) *User {
	if userIn, ok := r.people.Load(userId); ok {
		return userIn.(*User)
	}
	return nil
}

func (r *Room) GetRoomUsers() string {
	peopleByte, _ := json.Marshal(r.peopleList)
	return string(peopleByte)
}

func (r *Room) removeUser(userId string) {
	r.people.Delete(userId)
	for index, value := range r.peopleList {
		if value == userId {
			r.peopleList = append(r.peopleList[:index], r.peopleList[index+1:]...)
		}
	}
}

func (r *Room) AddUser(user *User) {
	r.people.Store(user.Id, user)
	r.peopleList = append(r.peopleList, user.Id)
}

func (r *Room) RunBattle() {
	log.Println(fmt.Sprintf("%s", "start battle"))
	peopleByte, _ := json.Marshal(r.peopleList)
	r.Broadcast(string(peopleByte), "startBattle")
	speed := StartSpeed
	for {
		select {
		// todo 对局结束判断
		default:
			if speed <= MaxSpeed {
				speed += SpeedStep
			}
			r.Broadcast(fmt.Sprintf("%f", speed), "speed")
		}
		time.Sleep(time.Second)
	}
}

func (r *Room) Broadcast(msgString, event string) {
	if msg.Message == nil {
		panic("未初始化广播引擎")
	}
	// 广播
	msg.Message.AddMsg(&msg.Msg{
		Msg:   msgString,
		Room:  r.Name,
		Event: event,
	})
}
