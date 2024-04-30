package Gpthread

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global/GpGame"
)

func Chatroom(Gthis *GpManager.WebSocketListController) {
	for {
		select {
		case JoinSocket := <-Gthis.SocketChan:
			if !Gthis.IsExistSocketById(JoinSocket.SocketId) {
				Gthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gthis, GpPacket.IM_S2C_JOIN, JoinSocket.User, JoinSocket.SocketId, "")

				Gthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("New socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			} else {
				Gthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("Old socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			}
		case SocketMessage := <-Gthis.MsgList:
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_EVENT_MESSAGE,
				GpPacket.IM_C2S2C_HEART: // 心跳
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, Gthis)
				break
			case
				GpPacket.IM_S2C_LEAVEROOM: // 关闭房间完成比赛，退出游戏，展示结果
				break
			case
				GpPacket.IM_S2C_JOIN: //创建房间，显示房间列表
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				break
			case
				GpPacket.IM_S2C_JOINCREATROOM: //创建房间等待对手加入
				break
			case
				GpPacket.IM_C2S_GETROOMLIST: // 刷新房间列
				GetRoomList(SocketMessage, Gthis)
				break
			case
				GpPacket.IM_EVENT_BROADCAST_HEART:
				Gphandle.GWebSocketStruct.BroadcastWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				break
			case
				GpPacket.IM_EVENT_BROADCAST_MESSAGE:
				Gphandle.GWebSocketStruct.BroadcastWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				break
			default:
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				break
			}
			GpPacket.NewArchive(SocketMessage)
			if SocketMessage.Type == GpPacket.IM_EVENT_MESSAGE {
				Global.Logger.Info("Message from", SocketMessage.Users.From, ";Msg:", SocketMessage.Msg)
			}

		case LeaveSocket := <-Gthis.UnSocketChan:
			for sub := Gthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(GpManager.SocketInfo).SocketId == LeaveSocket.SocketId {
					Gthis.ActiveSocketList.Remove(sub)
					// Clone connection.
					ws := sub.Value.(GpManager.SocketInfo).Conn
					if ws != nil {
						ws.Close()
						Global.Logger.Error("WebSocket closed:", LeaveSocket)
					}

					Gthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gthis, GpPacket.IM_S2C_LEAVEROOM, sub.Value.(GpManager.SocketInfo).User, LeaveSocket.SocketId, "") // Publish a LEAVE event.
					//Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
					break
				}
			}
		}
	}
}

func GameStart(event GpPacket.IM_protocol, Gthis *GpManager.WebSocketListController) {
	if "" == event.Msg {
		return
	}
	ret, t := GpGame.Start(event)
	if true == t {
		BCGame(ret, Gthis)
	}

}
func GetRoomList(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	ret := GpGame.GetRoomList()
	event := GpPacket.Protocol_getroomlist{}
	event.Type = GpPacket.IM_S2C_SENDROOMLIST
	event.List = ret
	event.SocketId = event.SocketId
	data, err := json.Marshal(msg)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(GpManager.SocketInfo).Conn
		if ws != nil {

			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				Global.Logger.Trace("disconnected user:", sub.Value.(GpManager.SocketInfo).User)
				Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}

			}
		}
	}
}
