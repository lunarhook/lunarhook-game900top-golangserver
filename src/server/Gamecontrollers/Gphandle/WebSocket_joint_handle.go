package Gphandle

import (
	"encoding/json"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *Gamecontrollers.WebSocketListController) SocketJoin(SocketId uint32, ws *websocket.Conn) {
	if this.IsExistSocketById(SocketId) {
		for sub := this.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
			if sub.Value.(Gamecontrollers.SocketInfo).SocketId == SocketId {
				this.ActiveSocketList.Remove(sub)
				break
			}
		}
	}

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		NewSocketId := r.Uint32()
		if !this.IsExistSocketById(NewSocketId) {
			//这里就是整个用户存在的循环体积，先将用户放入订阅队列
			this.SocketChan <- Gamecontrollers.SocketInfo{NewSocketId, GpPacket.IM_protocol_user{}, ws}
			//预定函数结尾让用户离开， 因为有可能强行kick，所以有单独函数
			defer this.SocketLeave(NewSocketId)
			//停止NewSocketId获取
			break
		}
	}

	// 后续socket的所有消息都在这里执行，如果断开都走defer leave干掉用户，心跳也在这里，目前还不支持多窗口单一心跳，这个将来客户端修改，主要是nginx time out300秒
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		var info GpPacket.IM_protocol
		if err := json.Unmarshal([]byte(p), &info); err == nil {
			this.MsgList <- this.NewMsg(info.Type, info.Users, info.SocketId, string(info.Msg))
			//G.Logger.Info(info)

		} else {
			Global.Logger.Error("Join", err)
		}

	}

}
