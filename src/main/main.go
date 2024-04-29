package main

import (
	_ "fmt"
	"time"

	"github.com/astaxie/beego"
	Gamecontrollers "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

type MainController struct {
	beego.Controller
}

func main() {
	Global.Init_Logs()
	time.Sleep(2 * time.Second)
	//beego.Router("/tetris/", &MainController{},"get:Single")
	//beego.Router("/tetris/watch", &MainController{},"get:Watch")
	//beego.Router("/tetris/game", &MainController{},"get:Game")
	beego.Router("/", &Gamecontrollers.AppController{})
	beego.Router("/join", &Gamecontrollers.AppController{}, "post:Join")

	// WebSocket.
	beego.Router("/wss", &Gamecontrollers.WebSocketController{})
	beego.Router("/wss/socket", &Gamecontrollers.WebSocketController{}, "get:Socket")
	beego.Router("/ws", &Gamecontrollers.WebSocketController{})
	beego.Router("/ws/socket", &Gamecontrollers.WebSocketController{}, "get:Socket")
	beego.Run()

}
