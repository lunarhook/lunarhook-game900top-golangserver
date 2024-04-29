package Gphandle

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"time"
)

func (this *WebSocketStruct) NewMsg(Gpthis *GpManager.WebSocketListController, ep GpPacket.EventType, user GpPacket.IM_protocol_user, SocketId uint32, msg string) GpPacket.IM_protocol {
	return GpPacket.IM_protocol{ep, msg, SocketId, user, int(time.Now().Unix())}
}
