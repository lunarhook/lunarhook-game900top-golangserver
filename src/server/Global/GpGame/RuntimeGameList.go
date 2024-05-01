package GpGame

import (
	"encoding/json"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

type GameTopRoom struct {
	id       int      `json:"id"`
	wordlist []string `json:"wordlist"`
	playa    []string `json:"playa"`
	playb    []string `json:"playb"`
	scorceA  int      `json:"scorceA"`
	scorceB  int      `json:"scorceB"`
	onwait   bool
	runplay  bool
}

var gGameTop *([]GameTopRoom)

func GameTopRoom_tick() {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		if true == (*gGameTop)[i].runplay {

		}
	}
}
func BuildServerRoom() {
	var GameTopRoomList = (make([]GameTopRoom, 3))
	gGameTop = &GameTopRoomList
	lens := len(*gGameTop)
	var Roomid = 10000
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].runplay {
			(*gGameTop)[i].id = Roomid
			(*gGameTop)[i].scorceA = 0
			(*gGameTop)[i].scorceB = 0
			(*gGameTop)[i].runplay = false
			(*gGameTop)[i].onwait = false
			Roomid++
		}
	}

}

func GetRoomList() string {
	lens := len(*gGameTop)
	type reslist struct {
		List []int
	}
	rlist := reslist{}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].runplay || true == (*gGameTop)[i].onwait {
			rlist.List = append(rlist.List, (*gGameTop)[i].id)
		}
	}
	data, err := json.Marshal(rlist)
	if err != nil {
		Global.Logger.Error("GetRoomList:", err)
		return ""
	}
	return string(data)
}

func JoinCreatRoomById(id int) string {
	lens := len(*gGameTop)
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].runplay || id == (*gGameTop)[i].id || true == (*gGameTop)[i].onwait {
			return "JOIN"
		}
		if false == (*gGameTop)[i].runplay || id == (*gGameTop)[i].id || false == (*gGameTop)[i].onwait {
			return "CREATE"
		}
	}
	return ""
}
