package Game

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"math/rand"
	"time"

	GpPacket2 "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
)

type GameInit struct {
	Side       int    `json:"side"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Speed      int    `json:"speed"`
	Num_block  int    `json:"num_block"`
	Type_color string `json:"type_color"`
	//Ident           int      `json:"ident"`
	Direction       int      `json:"direction"`
	Grade           int      `json:"grade"`
	Over            bool     `json:"Over"`
	Arr_bX          []int    `json:"arr_bX"`
	Arr_bY          []int    `json:"arr_bY"`
	Arr_store_X     []int    `json:"arr_store_X"`
	Arr_store_Y     []int    `json:"arr_store_Y"`
	Arr_store_color []string `json:"arr_store_color"`
}

var gGame *GameInit
var loop bool

var MsgList = make(chan string, 100)
var qMsgList = make(chan string, 100)

func initBackground() {

}
func initgame(event GpPacket2.IM_protocol) {
	loop = true
	go GameRussia()
	time.Sleep(40 * time.Millisecond)
	gGame.Side = 35
	gGame.Width = 10
	gGame.Height = 15
	gGame.Speed = 400
	gGame.Type_color = "000000"
	//map black by client
	initBackground()
	initBlock()
	//开启游戏循环，尝试第一个tick

	field := make(map[string]interface{}, 0)
	jsons, error := json.Marshal(gGame)
	if error != nil {
		fmt.Println(error.Error())
	}
	json.Unmarshal([]byte(jsons), &field)
	field["action"] = "start"
	field["notify"] = "start"
	jsons, error = json.Marshal(field)
	MsgList <- string(jsons)
}

func Down_speed_up_tick() {
	flag_all_down := true
	flag_all_down = JudgeCollision_down()

	if flag_all_down {
		//gGame.initBackground()
		for i := 0; i < len(gGame.Arr_bY); i++ {
			gGame.Arr_bY[i] = gGame.Arr_bY[i] + 1
		}
	} else {
		for i := 0; i < len(gGame.Arr_bX); i++ {
			gGame.Arr_store_X = append(gGame.Arr_store_X, gGame.Arr_bX[i])
			gGame.Arr_store_Y = append(gGame.Arr_store_Y, gGame.Arr_bY[i])
			gGame.Arr_store_color = append(gGame.Arr_store_color, gGame.Type_color)
		}
		gGame.Arr_bX = gGame.Arr_bX[0:0]
		gGame.Arr_bY = gGame.Arr_bY[0:0]
		initBlock()
	}
	ClearUnderBlock()
	//gGame.drawBlock(this.Type_color)
	//gGame.drawStaticBlock()

	MsgReturn()
	gameover()
}
func MsgReturn() {
	field := make(map[string]interface{}, 0)
	jsons, error := json.Marshal(gGame)
	if error != nil {
		fmt.Println(error.Error())
	}
	json.Unmarshal([]byte(jsons), &field)
	field["action"] = "tick"
	field["notify"] = "run"
	jsons, error = json.Marshal(field)
	qMsgList <- string(jsons)
}
func initBlock() {
	createRandom("rColor") //生成颜色字符串，
	createRandom("rBlock")
}

func gameover() {
	for i := 0; i < len(gGame.Arr_store_X); i++ {
		if gGame.Arr_store_Y[i] == 0 {
			loop = false
			gGame.Over = true
			field := make(map[string]interface{}, 0)
			jsons, error := json.Marshal(gGame)
			if error != nil {
				fmt.Println(error.Error())
			}
			json.Unmarshal([]byte(jsons), &field)
			field["action"] = "tick"
			field["notify"] = "end"
			jsons, error = json.Marshal(field)
			MsgList <- string(jsons)
		}
	}
}
func up_change_direction() {
	if gGame.Num_block == 5 {
		return
	}

	arr_tempX := []int{}
	arr_tempY := []int{}
	//因为不知道是否能够变形成功，所以先存储起来
	for i := 0; i < len(gGame.Arr_bX); i++ {
		arr_tempX = append(arr_tempX, gGame.Arr_bX[i])
		arr_tempY = append(arr_tempY, gGame.Arr_bY[i])
	}
	gGame.Direction++
	//将中心坐标提取出来，变形都以当前中心为准
	var ax_temp int
	var ay_temp int
	ax_temp = gGame.Arr_bX[0]
	ay_temp = gGame.Arr_bY[0]

	gGame.Arr_bX = gGame.Arr_bX[0:0] //将数组清空
	gGame.Arr_bY = gGame.Arr_bY[0:0]

	if gGame.Num_block == 1 {

		switch gGame.Direction % 4 {
		case 1:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp+1, ay_temp+1, ay_temp+1)
			break
		case 2:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp, ax_temp)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp-1, ay_temp+1)
			break
		case 3:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp+1, ay_temp)
			break
		case 0:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp+1, ay_temp)
			break
		}
	}
	if gGame.Num_block == 2 {

		switch gGame.Direction % 4 {
		case 1:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp-1, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp-1, ay_temp)
			break
		case 2:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp, ax_temp-1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp+1, ay_temp+1)
			break
		case 3:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp+1, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp, ay_temp+1)
			break
		case 0:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp+1, ay_temp-1)
			break
		}
	}
	if gGame.Num_block == 3 {

		switch gGame.Direction % 4 {
		case 1:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp+1, ax_temp+2)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp, ay_temp)
			break
		case 2:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp, ax_temp)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp+1, ay_temp+2)
			break
		case 3:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp+1, ax_temp+2)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp, ay_temp)
			break
		case 0:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp, ax_temp)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp+1, ay_temp+2)
			break
		}
	}
	if gGame.Num_block == 4 {

		switch gGame.Direction % 4 {
		case 1:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp-1, ax_temp, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp, ay_temp+1, ay_temp+1)
			break
		case 2:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp+1, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp+1, ay_temp, ay_temp-1)
			break
		case 3:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp-1, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp, ay_temp-1)
			break
		case 0:
			gGame.Arr_bX = append(gGame.Arr_bX, ax_temp, ax_temp, ax_temp+1, ax_temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, ay_temp, ay_temp-1, ay_temp, ay_temp+1)
			break
		}
	}

	if !(JudgeCollision_other(-1) && JudgeCollision_down() && JudgeCollision_other(1)) { //如果变形不成功则执行下面代码
		gGame.Arr_bX = gGame.Arr_bX[0:0]
		gGame.Arr_bY = gGame.Arr_bY[0:0]
		for i := 0; i < len(arr_tempX); i++ {
			gGame.Arr_bX = append(gGame.Arr_bX, arr_tempX[i])
			gGame.Arr_bY = append(gGame.Arr_bY, arr_tempY[i])
		}
	}
	MsgReturn()
	//this.drawStaticBlock()
}

// 方向键为左右的左移动函数
func move(dir_temp int) {
	//initBackground()

	if dir_temp == 1 { //右
		flag_all_right := true
		flag_all_right = JudgeCollision_other(1)

		if flag_all_right {
			for i := 0; i < len(gGame.Arr_bY); i++ {
				gGame.Arr_bX[i] = gGame.Arr_bX[i] + 1
			}
		}
	} else {
		flag_all_left := true
		flag_all_left = JudgeCollision_other(-1)

		if flag_all_left {
			for i := 0; i < len(gGame.Arr_bY); i++ {
				gGame.Arr_bX[i] = gGame.Arr_bX[i] - 1
			}
		}
	}
	//drawBlock(gGame.Type_color)
	//drawStaticBlock()
	MsgReturn()
}

func createRandom(stype string) {
	temp := gGame.Width/2 - 1

	if stype == "rBlock" {
		gGame.Num_block = rand.Intn(5) + 1
		switch gGame.Num_block {
		case 1:
			gGame.Arr_bX = append(gGame.Arr_bX, temp, temp-1, temp, temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, 0, 1, 1, 1)
			break
		case 2:
			gGame.Arr_bX = append(gGame.Arr_bX, temp, temp-1, temp-1, temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, 1, 0, 1, 1)
			break
		case 3:
			gGame.Arr_bX = append(gGame.Arr_bX, temp, temp-1, temp+1, temp+2)
			gGame.Arr_bY = append(gGame.Arr_bY, 0, 0, 0, 0)
			break
		case 4:
			gGame.Arr_bX = append(gGame.Arr_bX, temp, temp-1, temp, temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, 0, 0, 1, 1)
			break
		case 5:
			gGame.Arr_bX = append(gGame.Arr_bX, temp, temp+1, temp, temp+1)
			gGame.Arr_bY = append(gGame.Arr_bY, 0, 0, 1, 1)
			break
		}
	}
	if stype == "rColor" {
		num_color := rand.Intn(8) + 1

		switch num_color {
		case 1:
			gGame.Type_color = "#3EF72A"
			break
		case 2:
			gGame.Type_color = "yellow"
			break
		case 3:
			gGame.Type_color = "#2FE0BF"
			break
		case 4:
			gGame.Type_color = "red"
			break
		case 5:
			gGame.Type_color = "gray"
			break
		case 6:
			gGame.Type_color = "#C932C6"
			break
		case 7:
			gGame.Type_color = "#FC751B"
			break
		case 8:
			gGame.Type_color = "#6E6EDD"
			break
		case 9:
			gGame.Type_color = "#F4E9E1"
			break
		}
	}
}

func JudgeCollision_down() bool {
	for i := 0; i < len(gGame.Arr_bX); i++ {
		if gGame.Arr_bY[i]+1 == gGame.Height {
			return false
		}
		if len(gGame.Arr_store_X) != 0 {
			for j := 0; j < len(gGame.Arr_store_X); j++ {
				if gGame.Arr_bX[i] == gGame.Arr_store_X[j] {
					if gGame.Arr_bY[i]+1 == gGame.Arr_store_Y[j] {
						return false
					}
				}

			}
		}
	}
	return true
}
func ClearUnderBlock() {
	//删除低层方块
	var arr_row []int
	var line_num int
	if len(gGame.Arr_store_X) != 0 {
		for j := gGame.Height - 1; j >= 0; j-- {
			for i := 0; i < len(gGame.Arr_store_color); i++ {
				if gGame.Arr_store_Y[i] == j {
					arr_row = append(arr_row, i)
				}
			}
			if len(arr_row) == gGame.Width {
				line_num = j
				break
			} else {
				arr_row = arr_row[0:0]
			}
		}
	}
	if len(arr_row) == gGame.Width {
		//计算成绩grade
		gGame.Grade++

		for i := 0; i < len(arr_row); i++ {
			gGame.Arr_store_X = append(gGame.Arr_store_X[:arr_row[i]-i], gGame.Arr_store_X[arr_row[i]-i+1:]...)
			gGame.Arr_store_Y = append(gGame.Arr_store_Y[:arr_row[i]-i], gGame.Arr_store_Y[arr_row[i]-i+1:]...)
			gGame.Arr_store_color = append(gGame.Arr_store_color[:arr_row[i]-i], gGame.Arr_store_color[arr_row[i]-i+1:]...)
		}

		//让上面的方块往下掉一格
		for i := 0; i < len(gGame.Arr_store_color); i++ {
			if gGame.Arr_store_Y[i] < line_num {
				gGame.Arr_store_Y[i] = gGame.Arr_store_Y[i] + 1
			}
		}
	}
}

func JudgeCollision_other(num int) bool {
	for i := 0; i < len(gGame.Arr_bX); i++ {
		if num == 1 {
			if gGame.Arr_bX[i] == gGame.Width-1 {
				return false
			}

		}
		if num == -1 {
			if gGame.Arr_bX[i] == 0 {
				return false
			}

		}
		if len(gGame.Arr_store_X) != 0 {
			for j := 0; j < len(gGame.Arr_store_X); j++ {
				if gGame.Arr_bY[i] == gGame.Arr_store_Y[j] {
					if gGame.Arr_bX[i]+num == gGame.Arr_store_X[j] {
						return false
					}
				}
			}
		}
	}
	return true
}

func Start(event GpPacket2.IM_protocol) (GpPacket2.IM_protocol, bool) {
	if false == loop && "start" == event.Msg {
		initgame(event)
	}
	if string("left") == event.Msg {
		move(1)
	}
	if string("right") == event.Msg {
		move(-1)
	}
	if string("change") == event.Msg {
		up_change_direction()
	}
	if true == loop || len(MsgList) > 0 {
		select {
		case i := <-qMsgList:
			event.Msg = i
			return event, true
			break
		case i := <-MsgList:
			event.Msg = i
			return event, true
		}
	}

	//这里返回要客户端重新开始游戏
	return event, false
}

func GameRussia() {
	gGame = &GameInit{}
	for {
		time.Sleep(1000 * time.Millisecond)
		if true == loop {
			Down_speed_up_tick()
		} else {
			return
		}
	}
}
