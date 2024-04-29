package Gphandle

import (
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *WebSocketStruct) SocketLeave(SocketId uint32, Gpthis *GpManager.WebSocketListController) {
	Gpthis.UnSocketChan <- GpManager.UnSocketId{SocketId}
	Global.Logger.Info("Socket Leave:", SocketId)
}
