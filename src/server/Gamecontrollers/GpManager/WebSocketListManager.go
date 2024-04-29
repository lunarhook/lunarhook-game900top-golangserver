package GpManager

import (
	"container/list"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
)

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

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

func (this *WebSocketListController) IsExistSocketById(SocketId uint32) bool {
	for sub := this.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(SocketInfo).SocketId == SocketId {
			return true
		}
	}
	return false
}
