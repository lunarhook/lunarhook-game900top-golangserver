package main

import (
	_ "fmt"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
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
	Gamecontrollers.Init()
	IGphandle.Init()
	// WebSocket.ssl
	beego.Router("/wss", IGphandle.GHandlemanager)
	beego.Router("/wss/socket", IGphandle.GHandlemanager, "get:Api_socket")
	// WebSocket
	beego.Router("/ws", IGphandle.GHandlemanager)
	beego.Router("/ws/socket", IGphandle.GHandlemanager, "get:Api_socket")
	beego.Run()

}
