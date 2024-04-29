package Game

import "time"

type GameTopRoom struct {
	Total    int      `json:"total"`
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
	var GameTopRoomList = make([]GameTopRoom, 100)
	gGameTop = &GameTopRoomList
	for {
		time.Sleep(1000 * time.Millisecond)
		if true == loop {

		} else {
			gGameTop = nil
			return
		}
	}
}
