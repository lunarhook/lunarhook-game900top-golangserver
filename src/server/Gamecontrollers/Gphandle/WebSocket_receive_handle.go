package Gphandle

import (
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"net/http"

	"github.com/gorilla/websocket"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketStruct) HandleSocket(Gpthis *GpManager.WebSocketListController) {

	SocketId, _ := Gpthis.GetUint32("SocketId")

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// Upgrade from http request to WebSocket.
	ws, err := upgrader.Upgrade(Gpthis.Ctx.ResponseWriter, Gpthis.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(Gpthis.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		Global.Logger.Error("Cannot setup deivce WebSocket connection:", err)
		return
	}
	// Join chat room. 后续所有的通信都不会在走这里而是走到join函数里循环
	this.SocketJoin(SocketId, ws, Gpthis)

}
