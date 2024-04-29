package main

import (
	_ "fmt"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
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
	beego.Router("/wss", &Gphandle.WebSocketController{})
	beego.Router("/wss/socket", &Gphandle.WebSocketController{}, "get:Socket")
	beego.Router("/ws", &Gphandle.WebSocketController{})
	beego.Router("/ws/socket", &Gphandle.WebSocketController{}, "get:Socket")
	beego.Run()

}
