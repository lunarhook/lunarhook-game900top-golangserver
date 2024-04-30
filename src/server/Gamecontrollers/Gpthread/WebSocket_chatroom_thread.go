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

func Chatroom(Gpthis *GpManager.WebSocketListController) {
	for {
		select {
		case JoinSocket := <-Gpthis.SocketChan:
			if !Gpthis.IsExistSocketById(JoinSocket.SocketId) {
				Gpthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_S2C_JOIN, JoinSocket.User, JoinSocket.SocketId, "")

				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("New socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			} else {
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("Old socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			}
		case SocketMessage := <-Gpthis.MsgList:
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_EVENT_MESSAGE,
				GpPacket.IM_C2S2C_HEART: // 心跳
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_S2C_LEAVEROOM: // 关闭房间完成比赛，退出游戏，展示结果
				break
			case
				GpPacket.IM_S2C_JOIN: //创建房间，显示房间列表
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				GetRoomList(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_C2S_JOINCREATROOM:
				Global.Logger.Info("IM_S2C_JOINCREATROOM:", SocketMessage)
				break
			case
				GpPacket.IM_S2C_JOINCREATROOM: //创建房间等待对手加入
				Global.Logger.Info("IM_S2C_JOINCREATROOM:", SocketMessage)
				break
			case
				GpPacket.IM_C2S_GETROOMLIST: // 刷新房间列
				GetRoomList(SocketMessage, Gpthis)
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

		case LeaveSocket := <-Gpthis.UnSocketChan:
			for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(GpManager.SocketInfo).SocketId == LeaveSocket.SocketId {
					Gpthis.ActiveSocketList.Remove(sub)
					// Clone connection.
					ws := sub.Value.(GpManager.SocketInfo).Conn
					if ws != nil {
						ws.Close()
						Global.Logger.Error("WebSocket closed:", LeaveSocket)
					}

					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_S2C_LEAVEROOM, sub.Value.(GpManager.SocketInfo).User, LeaveSocket.SocketId, "") // Publish a LEAVE event.
					//Gpthis.MsgList <- Gpthis.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
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
func joincreatroom(event GpPacket.IM_protocol, Gthis *GpManager.WebSocketListController) {

}
func GetRoomList(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	ret := GpGame.GetRoomList()
	event := GpPacket.Protocol_getroomlist{}
	event.Type = GpPacket.IM_S2C_SENDROOMLIST
	event.List = ret
	event.SocketId = msg.SocketId
	data, err := json.Marshal(event)
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
