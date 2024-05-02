package Gphandle

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *WebSocketStruct) HeartWebSocket(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_C2S2C_HEART
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}

	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
			ws := sub.Value.(GpManager.SocketInfo).Conn
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// socket disconnected.
					Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
				} else {
					Global.Logger.Trace("S2C:", event)
				}
			}
		}
	}
}
