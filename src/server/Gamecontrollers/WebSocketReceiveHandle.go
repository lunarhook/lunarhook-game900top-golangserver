package Gamecontrollers

import (
	"net/http"

	"github.com/gorilla/websocket"
	Global2 "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Socket() {

	SocketId, _ := this.GetUint32("SocketId")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// Upgrade from http request to WebSocket.
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		Global2.Logger.Error("Cannot setup deivce WebSocket connection:", err)
		return
	}
	// Join chat room. 后续所有的通信都不会在走这里而是走到join函数里循环
	globaWebSocketListManager.SocketJoin(SocketId, ws)

}
