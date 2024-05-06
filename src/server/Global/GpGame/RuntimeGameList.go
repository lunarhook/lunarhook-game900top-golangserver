package GpGame

import (
	"encoding/json"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
	"github.com/tidwall/gjson"
	"time"
)

type GameTopRoom struct {
	Id          uint32    `json:"Id"`
	Wordlist    [5]string `json:"Wordlist"`
	SocketIdA   uint32    `json:"SocketIdA"`
	SocketIdB   uint32    `json:"SocketIdB"`
	PlayAavatar string    `json:"PlayAavatar"`
	PlayBavatar string    `json:"PlayBavatar"`
	ScorceA     uint32    `json:"ScorceA"`
	ScorceB     uint32    `json:"ScorceB"`
	TimeOut     uint32    `json:"TimeOut"`
	Runplay     bool
	Count       uint32
}

type BeginGameRoom struct {
	Id        uint32 `json:"Id"`
	SocketIdA uint32 `json:"SocketIdA"`
	SocketIdB uint32 `json:"SocketIdB"`
}

var gGameTop *([]GameTopRoom)

func GameTopRoom_tick() {

	for {
		time.Sleep(1 * time.Second)
		lens := len(*gGameTop)
		for i := 0; i < lens; i++ {
			if false == (*gGameTop)[i].Runplay {
			} else {
				Gpthis := GpManager.GlobaWebSocketListManager
				pGameRoom := &((*gGameTop)[i])
				pGameRoom.TimeOut = pGameRoom.TimeOut - 1
				//结尾是0的时候做下一个动作
				next := (pGameRoom.TimeOut%10 == 0)
				//检查是否发词
				if next == true && pGameRoom.TimeOut > 0 {
					step := (50 - pGameRoom.TimeOut) / 10
					wordlist := pGameRoom.Wordlist[step]
					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_SENDQUEST, (*gGameTop)[i].SocketIdA, string(wordlist), pGameRoom.Id)
					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_SENDQUEST, (*gGameTop)[i].SocketIdB, string(wordlist), pGameRoom.Id)
				}
				//回收房间
				if pGameRoom.TimeOut <= 0 || pGameRoom.Count >= 2 {
					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_LEAVEROOM, (*gGameTop)[i].SocketIdA, ("RoomClose!"), 0)
					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_LEAVEROOM, (*gGameTop)[i].SocketIdB, ("RoomClose!"), 0)
					Clearroom(pGameRoom)
				}
				Global.Logger.Info("S2S: GameTopRoom_tick =", pGameRoom.Id, pGameRoom.TimeOut)
			}
		}
	}

}
func Clearroom(pGroom *GameTopRoom) {
	pGroom.TimeOut = 0
	pGroom.Runplay = false
	pGroom.SocketIdA = 0
	pGroom.SocketIdB = 0
	pGroom.ScorceA = 0
	pGroom.ScorceB = 0
	pGroom.Wordlist = [5]string{}
	pGroom.PlayAavatar = ""
	pGroom.PlayBavatar = ""
	pGroom.Count = 0
}
func BuildServerRoom(size uint) {
	var GameTopRoomList = (make([]GameTopRoom, size))
	gGameTop = &GameTopRoomList
	lens := len(*gGameTop)
	var Roomid = uint32(10000)
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
func ActiveRoom(msg GpPacket.IM_rec) {
	lens := len(*gGameTop)
	if gjson.Valid(msg.Msg) {
		result := gjson.Get(msg.Msg, "Msg.quest")
		Global.Logger.Info("ActiveRoom:", result)
		wordlist := [5]string{}
		for i, name := range result.Array() {
			wordlist[i] = name.String()
		}

		for i := 0; i < lens; i++ {
			if msg.SocketId == (*gGameTop)[i].SocketIdA {
				(*gGameTop)[i].Runplay = true
				(*gGameTop)[i].TimeOut = 51
				(*gGameTop)[i].Wordlist = wordlist
				break
			} else if msg.SocketId == (*gGameTop)[i].SocketIdB {
				(*gGameTop)[i].Runplay = true
				(*gGameTop)[i].TimeOut = 51
				(*gGameTop)[i].Wordlist = wordlist
				break
			} else {

			}
		}

	}

}
func GetRoomList() string {
	lens := len(*gGameTop)
	type reslist struct {
		RoomId   uint32 `json:"RoomId"`
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
func JoinCreatRoomById(id uint32, mysock uint32) *GameTopRoom {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		//join
		if (*gGameTop)[i].SocketIdA == mysock {
			return &(*gGameTop)[i]
		} else if (*gGameTop)[i].SocketIdB == mysock {
			return &(*gGameTop)[i]
		}
	}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].Runplay && id == (*gGameTop)[i].Id {
			return &(*gGameTop)[i]
		}
	}
	return nil
}
func GetRoomById(id uint32, mysock uint32) *GameTopRoom {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		if (*gGameTop)[i].SocketIdA == mysock || (*gGameTop)[i].SocketIdB == mysock {
			return &(*gGameTop)[i]
		}
	}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].Runplay && id == (*gGameTop)[i].Id {
			return &(*gGameTop)[i]
		}
	}
	return nil
}
