package Gphandle

import (
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func SocketLeave(SocketId uint32, Gpthis *Gamecontrollers.WebSocketListController) {
	Gpthis.UnSocketChan <- Gamecontrollers.UnSocketId{SocketId}
	Global.Logger.Info("Socket Leave:", SocketId)
}
