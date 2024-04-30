package GpGame

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
	ClearUnderBlock()
	MsgReturn()
	gameover()
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

func GetRoomList() []int {
	lens := len(*gGameTop)
	reslist := []int{}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].runplay || true == (*gGameTop)[i].onwait {
			reslist = append(reslist, (*gGameTop)[i].id)
		}
	}
	return reslist
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
