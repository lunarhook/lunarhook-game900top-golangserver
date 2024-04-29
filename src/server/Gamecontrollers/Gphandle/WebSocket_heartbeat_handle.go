package Gphandle

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func HeartWebSocket(event GpPacket.IM_protocol, Gpthis *Gamecontrollers.WebSocketListController) {
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}

	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		if sub.Value.(Gamecontrollers.SocketInfo).SocketId == event.SocketId {
			ws := sub.Value.(Gamecontrollers.SocketInfo).Conn
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// socket disconnected.
					Gpthis.UnSocketChan <- Gamecontrollers.UnSocketId{sub.Value.(Gamecontrollers.SocketInfo).SocketId}
				} else {
					Global.Logger.Trace("Socketheart :", event.SocketId)
				}
			}
		}
	}
}