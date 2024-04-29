package Gamecontrollers

import (
	"container/list"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gpthread"
)

func Init() {
	Gphandle.GWebSocketStruct = &Gphandle.WebSocketStruct{}
	GpManager.GlobaWebSocketListManager = &GpManager.WebSocketListController{}
	GpManager.GlobaWebSocketListManager.SocketChan = make(chan GpManager.SocketInfo, 100)
	// Channel for exit users.
	GpManager.GlobaWebSocketListManager.UnSocketChan = make(chan GpManager.UnSocketId, 100)
	// Send events here to publish them.
	GpManager.GlobaWebSocketListManager.MsgList = make(chan GpPacket.IM_protocol, 100)

	GpManager.GlobaWebSocketListManager.ActiveSocketList = list.New()

	go Gpthread.Chatroom(GpManager.GlobaWebSocketListManager)
	go Gpthread.NetRussia(GpManager.GlobaWebSocketListManager)
}
