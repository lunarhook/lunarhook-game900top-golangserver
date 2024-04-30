package GpGame

type GameTopRoom struct {
	id       int      `json:"id"`
	wordlist []string `json:"wordlist"`
	playa    []string `json:"playa"`
	playb    []string `json:"playb"`
	scorceA  int      `json:"scorceA"`
	scorceB  int      `json:"scorceB"`
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
}

func GetRoomList() []int {
	lens := len(*gGameTop)
	reslist := []int{}
	for i := 0; i < lens; i++ {
		if false == (*gGameTop)[i].runplay {
			reslist = append(reslist, (*gGameTop)[i].id)
		}
	}
	return reslist
}
