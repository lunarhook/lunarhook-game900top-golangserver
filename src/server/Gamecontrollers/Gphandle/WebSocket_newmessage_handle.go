package Gphandle

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"time"
)

func (this *WebSocketStruct) NewMsg(Gpthis *GpManager.WebSocketListController, ep GpPacket.EventType, user GpPacket.IM_protocol, SocketId uint32, msg string) GpPacket.IM_protocol {
	return GpPacket.IM_protocol{ep, msg, SocketId, 0, user.Users, int(time.Now().Unix()), 0}
}
func (this *WebSocketStruct) NewByte(Gpthis *GpManager.WebSocketListController, ep GpPacket.EventType, b []byte) GpPacket.IM_msg {
	return GpPacket.IM_msg{ep, b}
}
