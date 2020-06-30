package match

import (
	"fmt"
	"log"
	"sync"
	"time"
	"trex_online/src/util"
)

const (
	RoomId int32 = 0
)

type Room struct {
	Name string // 房间名称
	Id   int32

	Cap    int32    // 容量
	Len    int32    // 数量
	Finish bool     // 对局是否完成
	Full   bool     // 是否满员
	People sync.Map // 对局的用户

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
		People:     sync.Map{},
		CreateTime: time.Now().Unix(),
		FinishTime: 0,
	}
	return room
}

func (r *Room) RunBattle() {
	log.Println(fmt.Sprintf("%s", "start battle"))
	return
}

func (r *Room) Broadcast() {
	// todo 广播

}
