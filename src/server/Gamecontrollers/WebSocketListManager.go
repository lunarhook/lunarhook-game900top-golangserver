package Gamecontrollers

import (
	"container/list"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gpthread"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global/Game"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

var GlobaWebSocketListManager *WebSocketListController

// 带用户信息的websocket
type SocketInfo struct {
	SocketId uint32
	User     GpPacket.IM_protocol_user
	Conn     *websocket.Conn
}
type SocketId struct {
	SocketId uint32
}
type UnSocketId struct {
	SocketId uint32
}
type WebSocketListController struct {
	// Channel for new join users.
	SocketChan chan SocketInfo
	// Channel for exit users.
	UnSocketChan chan UnSocketId
	// Send events here to publish them.
	MsgList chan (GpPacket.IM_protocol)
	// Long polling waiting list.
	ActiveSocketList *list.List
	beego.Controller
}

func (this *WebSocketListController) IsExistSocketById(SocketId uint32) bool {
	for sub := this.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(SocketInfo).SocketId == SocketId {
			return true
		}
	}
	return false
}

func init() {
	GlobaWebSocketListManager = &WebSocketListController{}
	GlobaWebSocketListManager.SocketChan = make(chan SocketInfo, 100)
	// Channel for exit users.
	GlobaWebSocketListManager.UnSocketChan = make(chan UnSocketId, 100)
	// Send events here to publish them.
	GlobaWebSocketListManager.MsgList = make(chan GpPacket.IM_protocol, 100)

	GlobaWebSocketListManager.ActiveSocketList = list.New()

	go Gpthread.Chatroom(GlobaWebSocketListManager)
	go Gpthread.NetRussia(GlobaWebSocketListManager)
}

func (this *WebSocketListController) NewMsg(ep GpPacket.EventType, user GpPacket.IM_protocol_user, SocketId uint32, msg string) GpPacket.IM_protocol {
	return GpPacket.IM_protocol{ep, msg, SocketId, user, int(time.Now().Unix())}
}

func (this *WebSocketListController) BCGame(event GpPacket.IM_protocol) {
	Gphandle.BroadcastWebSocket(event, GlobaWebSocketListManager)
}

func (this *WebSocketListController) Game(event GpPacket.IM_protocol) {
	if "" == event.Msg {
		return
	}
	ret, t := Game.Start(event)
	if true == t {
		this.BCGame(ret)
	}

}
