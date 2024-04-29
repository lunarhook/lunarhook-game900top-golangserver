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

func (ih *Handlemanager) Api_socket() {
	ih.ph.HandleSocket(ih.pC)
}
