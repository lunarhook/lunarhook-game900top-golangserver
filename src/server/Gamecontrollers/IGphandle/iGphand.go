package IGphandle

import (
	"github.com/astaxie/beego"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
)

type Handlemanager struct {
	pC *GpManager.WebSocketListController
	ph *Gphandle.WebSocketStruct
	beego.Controller
}

func (p *Handlemanager) Api_socket() {
	GHandlemanager.Controller = p.Controller
	GHandlemanager.pC.Controller = p.Controller
	GHandlemanager.ph.Controller = p.Controller
	GHandlemanager.ph.HandleSocket(GHandlemanager.pC)
}

var GHandlemanager *Handlemanager

func Init() {
	GHandlemanager = &Handlemanager{}
	GHandlemanager.pC = GpManager.GlobaWebSocketListManager
	GHandlemanager.ph = Gphandle.GWebSocketStruct
}
