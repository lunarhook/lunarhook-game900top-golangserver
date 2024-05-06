package Gphandle

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
)

func (this *WebSocketStruct) NewByte(Gpthis *GpManager.WebSocketListController, ep GpPacket.EventType, SocketId uint32, msg string, rid uint32) GpPacket.IM_rec {
	return GpPacket.IM_rec{ep, SocketId, msg, rid}
}
