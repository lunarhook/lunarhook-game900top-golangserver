package Gamecontrollers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *WebSocketListController) HeartWebSocket(event GpPacket.IM_protocol) {
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}

	for sub := this.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		if sub.Value.(SocketInfo).SocketId == event.SocketId {
			ws := sub.Value.(SocketInfo).Conn
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// socket disconnected.
					this.UnSocketChan <- UnSocketId{sub.Value.(SocketInfo).SocketId}
				} else {
					Global.Logger.Trace("Socketheart :", event.SocketId)
				}
			}
		}
	}
}
