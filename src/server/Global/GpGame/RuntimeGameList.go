package GpGame

import (
	"encoding/json"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
	"time"
)

type GameTopRoom struct {
	Id          uint64   `json:"Id"`
	Wordlist    []string `json:"Wordlist"`
	SocketIdA   uint32   `json:"SocketIdA"`
	SocketIdB   uint32   `json:"SocketIdB"`
	PlayAavatar string   `json:"PlayAavatar"`
	PlayBavatar string   `json:"PlayBavatar"`
	ScorceA     uint64   `json:"ScorceA"`
	ScorceB     uint64   `json:"ScorceB"`
	Runplay     bool
}

type BeginGameRoom struct {
	Id        uint64 `json:"Id"`
	SocketIdA uint32 `json:"SocketIdA"`
	SocketIdB uint32 `json:"SocketIdB"`
}

var gGameTop *([]GameTopRoom)

func GameTopRoom_tick() {

	for {
		time.Sleep(10 * time.Second)
		lens := len(*gGameTop)
		for i := 0; i < lens; i++ {
			if false == (*gGameTop)[i].Runplay {
				if 0 != (*gGameTop)[i].SocketIdA && (*gGameTop)[i].SocketIdB == 0 {
					Gpthis := GpManager.GlobaWebSocketListManager
					o := BeginGameRoom{}
					o.Id = (*gGameTop)[i].Id
					o.SocketIdA = (*gGameTop)[i].SocketIdA
					o.SocketIdB = (*gGameTop)[i].SocketIdA
					data, _ := json.Marshal(o)
					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_C2S_TESTBEGINGAME, (*gGameTop)[i].SocketIdA, string(data))
				}
			}
		}
	}

}
func BuildServerRoom(size uint) {
	var GameTopRoomList = (make([]GameTopRoom, size))
	gGameTop = &GameTopRoomList
	lens := len(*gGameTop)
	var Roomid = uint64(10000)
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].Runplay {
			(*gGameTop)[i].Id = Roomid
			(*gGameTop)[i].ScorceA = 0
			(*gGameTop)[i].ScorceB = 0
			(*gGameTop)[i].Runplay = false
			Roomid++
		}
	}

}

func GetRoomList() string {
	lens := len(*gGameTop)
	type reslist struct {
		RoomId   uint64 `json:"RoomId"`
		Battle   bool   `json:"Battle"`
		SocketId uint32 `json:"SocketId"`
	}
	rlist := []reslist{}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].Runplay {
			index := reslist{}
			index.RoomId = (*gGameTop)[i].Id
			index.Battle = ((0 != (*gGameTop)[i].SocketIdA) || (0 != (*gGameTop)[i].SocketIdB))
			if 0 != (*gGameTop)[i].SocketIdA {
				index.SocketId = (*gGameTop)[i].SocketIdA
			} else {
				index.SocketId = (*gGameTop)[i].SocketIdB
			}
			rlist = append(rlist, index)
		}
	}
	data, err := json.Marshal(rlist)
	if err != nil {
		Global.Logger.Error("GetRoomList:", err)
		return ""
	}
	return string(data)
}
func LeaveRoomBySocketId(SocketId uint32) {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		if SocketId == (*gGameTop)[i].SocketIdA {
			(*gGameTop)[i].SocketIdA = 0
		}
		if SocketId == (*gGameTop)[i].SocketIdB {
			(*gGameTop)[i].SocketIdB = 0
		}
	}
}
func JoinCreatRoomById(id uint64, mysock uint32) *GameTopRoom {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		if (*gGameTop)[i].SocketIdA == mysock || (*gGameTop)[i].SocketIdB == mysock {
			return nil
		}
	}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].Runplay && id == (*gGameTop)[i].Id {
			return &(*gGameTop)[i]
		}
	}
	return nil
}
