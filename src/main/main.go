package main

import (
	_ "fmt"
	"time"

	"github.com/astaxie/beego"
	Gamecontrollers2 "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	Global2 "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

type MainController struct {
	beego.Controller
}

func main() {
	Global2.Init_Logs()
	time.Sleep(2 * time.Second)
	//beego.Router("/tetris/", &MainController{},"get:Single")
	//beego.Router("/tetris/watch", &MainController{},"get:Watch")
	//beego.Router("/tetris/game", &MainController{},"get:Game")
	beego.Router("/", &Gamecontrollers2.AppController{})
	beego.Router("/join", &Gamecontrollers2.AppController{}, "post:Join")

	// WebSocket.
	beego.Router("/wss", &Gamecontrollers2.WebSocketController{})
	beego.Router("/wss/socket", &Gamecontrollers2.WebSocketController{}, "get:Socket")
	beego.Router("/ws", &Gamecontrollers2.WebSocketController{})
	beego.Router("/ws/socket", &Gamecontrollers2.WebSocketController{}, "get:Socket")
	beego.Run()

}
