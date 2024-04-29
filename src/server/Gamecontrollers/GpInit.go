package Gamecontrollers

import (
	"container/list"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gpthread"
)

func init() {
	var GlobaWebSocketListManager *GpManager.WebSocketListController
	GlobaWebSocketListManager = &GpManager.WebSocketListController{}
	GlobaWebSocketListManager.SocketChan = make(chan GpManager.SocketInfo, 100)
	// Channel for exit users.
	GlobaWebSocketListManager.UnSocketChan = make(chan GpManager.UnSocketId, 100)
	// Send events here to publish them.
	GlobaWebSocketListManager.MsgList = make(chan GpPacket.IM_protocol, 100)

	GlobaWebSocketListManager.ActiveSocketList = list.New()

	go Gpthread.Chatroom(GlobaWebSocketListManager)
	go Gpthread.NetRussia(GlobaWebSocketListManager)
}
