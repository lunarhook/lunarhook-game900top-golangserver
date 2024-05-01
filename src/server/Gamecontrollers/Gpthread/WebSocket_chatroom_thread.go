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
			Global.Logger.Info("C2S:", JoinSocket)
			if !Gpthis.IsExistSocketById(JoinSocket.SocketId) {
				Gpthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_JOIN, JoinSocket.SocketId, string("first join"))
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.SocketId, string("welcome"))
			} else {
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.SocketId, string("welcome again"))
			}
		case SocketMessage := <-Gpthis.MsgList:
			Global.Logger.Info("C2S:", SocketMessage)
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
				Global.Logger.Info("S2C:", SocketMessage)
			}

		case LeaveSocket := <-Gpthis.UnSocketChan:
			Global.Logger.Info("S2C:", LeaveSocket)
			for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(GpManager.SocketInfo).SocketId == LeaveSocket.SocketId {
					Gpthis.ActiveSocketList.Remove(sub)
					// Clone connection.
					ws := sub.Value.(GpManager.SocketInfo).Conn
					if ws != nil {
						ws.Close()
						Global.Logger.Error("WebSocket closed:", LeaveSocket)
					}

					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_LEAVEROOM, LeaveSocket.SocketId, ("good bye!")) // Publish a LEAVE event.
					//Gpthis.MsgList <- Gpthis.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
					break
				}
			}
		}
	}
}

func SelectRoom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_ret{}
	event.Type = GpPacket.IM_S2C_JOINCREATROOM
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
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
func joinroom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_ret{}
	event.Type = GpPacket.IM_S2C_JOIN
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
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
func LeaveRoom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_ret{}
	event.Type = GpPacket.IM_S2C_LEAVEROOM
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
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
func GetRoomList(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	ret := GpGame.GetRoomList()
	event := GpPacket.IM_ret{}
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
