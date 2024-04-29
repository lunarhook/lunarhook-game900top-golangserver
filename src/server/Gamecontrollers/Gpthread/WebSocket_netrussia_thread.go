package Gpthread

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global/GpGame"
	"time"
)

func NetRussia(Gthis *Gamecontrollers.WebSocketListController) {
	for {
		time.Sleep(400 * time.Millisecond)
		event := GpPacket.IM_protocol{}
		event.Type = GpPacket.IM_EVENT_BROADCAST_MESSAGE
		ret, b := GpGame.Start(event)
		if true == b {
			BCGame(ret, Gthis)
		}
	}
}
func BCGame(event GpPacket.IM_protocol, Gthis *Gamecontrollers.WebSocketListController) {
	Gphandle.BroadcastWebSocket(event, Gthis)
}
