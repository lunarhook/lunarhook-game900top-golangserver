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
			Global.Logger.Info("Join:", JoinSocket)
			if !Gpthis.IsExistSocketById(JoinSocket.SocketId) {
				Gpthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_S2C_JOIN, JoinSocket.User, JoinSocket.SocketId, "first")

				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("New socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			} else {
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewMsg(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome again")
				Global.Logger.Info("Old socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			}
		case SocketMessage := <-Gpthis.MsgList:
			Global.Logger.Info("Msg:", SocketMessage)
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_C2S2C_HEART: // 心跳
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_C2S_LEAVEROOM,
				GpPacket.IM_S2C_LEAVEROOM: // 关闭房间完成比赛，退出游戏，展示结果
				LeaveRoom(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_S2C_JOIN: //创建房间，显示房间列表
				joinroom(SocketMessage, Gpthis)
				//GetRoomList(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_C2S_JOINCREATROOM:
				SelectRoom(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_S2C_JOINCREATROOM: //创建房间等待对手加入
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
				Global.Logger.Info("IM_EVENT_MESSAGE", SocketMessage)
			}

		case LeaveSocket := <-Gpthis.UnSocketChan:
			Global.Logger.Info("Leave:", LeaveSocket)
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

func GameStart(event GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	if "" == event.Msg {
		return
	}
	ret, t := GpGame.Start(event)
	if true == t {
		BCGame(ret, Gpthis)
	}

}
func SelectRoom(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_protocol{}
	event.Type = GpPacket.IM_S2C_JOINCREATROOM
	event.SocketId = msg.SocketId
	event.SocketIdother = msg.SocketIdother
	field := make(map[string]interface{}, 0)
	ret, err := json.Marshal(msg.Msg)
	err = json.Unmarshal(ret, field)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(GpManager.SocketInfo).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, []byte(ret)) != nil {
				Global.Logger.Trace("disconnected user:", sub.Value.(GpManager.SocketInfo).User)
				// User disconnected.
				Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
			}

		}
	}
}
func joinroom(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_protocol{}
	event.Type = GpPacket.IM_S2C_JOIN
	event.SocketId = msg.SocketId
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(GpManager.SocketInfo).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				Global.Logger.Trace("disconnected user:", sub.Value.(GpManager.SocketInfo).User)
				// User disconnected.
				Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
			}

		}
	}
}
func LeaveRoom(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_protocol{}
	event.Type = GpPacket.IM_S2C_LEAVEROOM
	event.SocketId = msg.SocketId
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		ws := sub.Value.(GpManager.SocketInfo).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				Global.Logger.Trace("disconnected user:", sub.Value.(GpManager.SocketInfo).User)
			}
			// User disconnected.
			Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
		}
	}
}
func GetRoomList(msg GpPacket.IM_protocol, Gpthis *GpManager.WebSocketListController) {
	ret := GpGame.GetRoomList()
	event := GpPacket.IM_protocol{}
	event.Type = GpPacket.IM_S2C_SENDROOMLIST
	event.Msg = ret
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
