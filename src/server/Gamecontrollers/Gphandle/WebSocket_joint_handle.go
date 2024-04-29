package Gphandle

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
	"math/rand"
	"time"
)

func (this *WebSocketStruct) SocketJoin(SocketId uint32, ws *websocket.Conn, Gpthis *GpManager.WebSocketListController) {
	if Gpthis.IsExistSocketById(SocketId) {
		for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
			if sub.Value.(GpManager.SocketInfo).SocketId == SocketId {
				Gpthis.ActiveSocketList.Remove(sub)
				break
			}
		}
	}

	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		NewSocketId := r.Uint32()
		if !Gpthis.IsExistSocketById(NewSocketId) {
			//这里就是整个用户存在的循环体积，先将用户放入订阅队列
			Gpthis.SocketChan <- GpManager.SocketInfo{NewSocketId, GpPacket.IM_protocol_user{}, ws}
			//预定函数结尾让用户离开， 因为有可能强行kick，所以有单独函数
			defer this.SocketLeave(NewSocketId, Gpthis)
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
			Gpthis.MsgList <- this.NewMsg(Gpthis, info.Type, info.Users, info.SocketId, string(info.Msg))
			//G.Logger.Info(info)

		} else {
			Global.Logger.Error("Join", err)
		}

	}

}
