package Gphandle

import (
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *Gamecontrollers.WebSocketListController) SocketLeave(SocketId uint32) {
	this.UnSocketChan <- Gamecontrollers.UnSocketId{SocketId}
	Global.Logger.Info("Socket Leave:", SocketId)
}
