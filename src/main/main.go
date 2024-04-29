package main

import (
	_ "fmt"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/IGphandle"
	"time"

	"github.com/astaxie/beego"
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
	beego.Router("/", &GpManager.AppController{})
	beego.Router("/join", &GpManager.AppController{}, "post:Join")

	// WebSocket.
	beego.Router("/wss", &GpManager.WebSocketController{})
	beego.Router("/wss/socket", &IGphandle.Handlemanager{}, "get:Api_socket")
	beego.Router("/ws", &GpManager.WebSocketController{})
	beego.Router("/ws/socket", &IGphandle.Handlemanager{}, "get:Api_socket")
	beego.Run()

}
