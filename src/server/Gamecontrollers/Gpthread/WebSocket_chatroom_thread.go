package Gpthread

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/GpManager"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global/GpGame"
	"github.com/tidwall/gjson"
)

func Chatroom(Gpthis *GpManager.WebSocketListController) {
	for {
		select {
		case JoinSocket := <-Gpthis.SocketChan:
			Global.Logger.Info("C2S:", JoinSocket)
			if !Gpthis.IsExistSocketById(JoinSocket.SocketId) {
				Gpthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_JOIN, JoinSocket.SocketId, string("first join"), 0)
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.SocketId, string("welcome"), 0)
			} else {
				Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_EVENT_MESSAGE, JoinSocket.SocketId, string("welcome again"), 0)
			}
		case SocketMessage := <-Gpthis.MsgList:
			Global.Logger.Info("C2S:", SocketMessage)
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_C2S2C_HEART: // 心跳
				Gphandle.GWebSocketStruct.HeartWebSocket(SocketMessage, Gpthis)
				GetRoomList(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_C2S_LEAVEROOM:
				closesocket(SocketMessage, Gpthis)
				break
			case
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
				GetRoomList(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_S2C_JOINCREATROOM: //创建房间等待对手加入
				break
			case
				GpPacket.IM_C2S_GETROOMLIST: // 刷新房间列
				GetRoomList(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_S2C_SENDQUEST:
				RoomGameRunNext(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_EVENT_BROADCAST_HEART:
				Gphandle.GWebSocketStruct.BroadcastWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				break
			case
				GpPacket.IM_C2S_SENDANSWER:
				RoomGameRunUpdate(SocketMessage, Gpthis)
				break
			case
				//广播所有消息事件
				GpPacket.IM_EVENT_BROADCAST_MESSAGE:
				Gphandle.GWebSocketStruct.BroadcastWebSocket(SocketMessage, GpManager.GlobaWebSocketListManager)
				GetRoomList(SocketMessage, Gpthis)
				break
			case
				//GpPacket.IM_C2S_TESTBEGINGAME,
				GpPacket.IM_S2C_BEGINQUEST:
				BeginRoomGame(SocketMessage, Gpthis)
				break
			case
				GpPacket.IM_C2S_SENDQUEST:
				RoomGameRun(SocketMessage, Gpthis)
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

					Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_LEAVEROOM, LeaveSocket.SocketId, ("good bye!"), 0) // Publish a LEAVE event.
					//Gpthis.MsgList <- Gpthis.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
					break
				}
			}
		}
	}
}

func SelectRoom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_JOINCREATROOM
	event.SocketId = msg.SocketId
	RoomId := gjson.Get(msg.Msg, "Msg.roomid")
	pRoom := GpGame.JoinCreatRoomById(uint32(RoomId.Uint()), msg.SocketId)
	if nil != pRoom {
		r := &(*pRoom)
		if 0 == r.SocketIdA {
			r.SocketIdA = msg.SocketId
			avatar := gjson.Get(msg.Msg, "Msg.avatarUrl")
			r.PlayAavatar = avatar.String()
		} else if 0 == r.SocketIdB {
			r.SocketIdB = msg.SocketId
			avatar := gjson.Get(msg.Msg, "Msg.avatarUrl")
			r.PlayBavatar = avatar.String()
		}
		if 0 != pRoom.SocketIdA && 0 != pRoom.SocketIdB {
			Gpthis := GpManager.GlobaWebSocketListManager
			o := &GpGame.BeginGameRoom{}
			o.Id = pRoom.Id
			o.SocketIdA = pRoom.SocketIdA
			o.SocketIdB = pRoom.SocketIdB
			//这里是让所有的数据都从房间的A玩家开始传递
			Gpthis.MsgList <- Gphandle.GWebSocketStruct.NewByte(Gpthis, GpPacket.IM_S2C_BEGINQUEST, pRoom.SocketIdA, string(""), pRoom.Id)
		}
		ret := map[string]interface{}{"AvatarA": pRoom.PlayAavatar, "AvatarB": pRoom.PlayBavatar, "SocketIdA": pRoom.SocketIdA, "SocketIdB": pRoom.SocketIdB}
		v, _ := json.Marshal(ret)
		event.Msg = string(v)
	} else {
		event.Msg = string("IM_S2C_JOINCREATROOM failed")
		return
	}
	data, err := json.Marshal(event)
	//JoinCreatRoomById
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == pRoom.SocketIdA || sub.Value.(GpManager.SocketInfo).SocketId == pRoom.SocketIdB {
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
}
func joinroom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_JOIN
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
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
}
func closesocket(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == msg.SocketId {
			ws := sub.Value.(GpManager.SocketInfo).Conn
			if ws != nil {
				Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
			}
		}
	}
}
func LeaveRoom(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_LEAVEROOM
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
	GpGame.LeaveRoomBySocketId(msg.SocketId)
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
			ws := sub.Value.(GpManager.SocketInfo).Conn
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					Global.Logger.Trace("disconnected user:", sub.Value.(GpManager.SocketInfo).User)
					Gpthis.UnSocketChan <- GpManager.UnSocketId{sub.Value.(GpManager.SocketInfo).SocketId}
				}

			}
		}
	}
}
func GetRoomList(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	ret := GpGame.GetRoomList()
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_SENDROOMLIST
	event.Msg = ret
	event.SocketId = msg.SocketId
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
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
}
func RoomGameRunUpdate(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_SENDANSWER
	event.SocketId = msg.SocketId
	event.RoomId = msg.RoomId
	pRoom := GpGame.GetRoomById(event.RoomId, event.SocketId)
	if pRoom.TimeOut < 10 {
		pRoom.Count = pRoom.Count + 1
	}
	answer := gjson.Get(msg.Msg, "Msg.answer").Str
	for i := 0; i < len(pRoom.Wordlist); i++ {
		ss := gjson.Get(pRoom.Wordlist[i], "mean").Str
		if ss == answer {
			if msg.SocketId == pRoom.SocketIdA {
				pRoom.ScorceA = pRoom.ScorceA + pRoom.TimeOut%10
			}
			if msg.SocketId == pRoom.SocketIdB {
				pRoom.ScorceB = pRoom.ScorceB + pRoom.TimeOut%10
			}
		}
	}
	ret := map[string]interface{}{"ScorceA": pRoom.ScorceA, "ScorceB": pRoom.ScorceB, "SocketIdA": pRoom.SocketIdA, "SocketIdB": pRoom.SocketIdB}
	v, _ := json.Marshal(ret)
	event.Msg = string(v)
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == pRoom.SocketIdA || sub.Value.(GpManager.SocketInfo).SocketId == pRoom.SocketIdB {
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
}
func RoomGameRunNext(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_SENDQUEST
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
	event.RoomId = msg.RoomId

	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
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
}
func RoomGameRun(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	GpGame.ActiveRoom(msg)
}
func BeginRoomGame(msg GpPacket.IM_rec, Gpthis *GpManager.WebSocketListController) {
	event := GpPacket.IM_rec{}
	event.Type = GpPacket.IM_S2C_BEGINQUEST
	event.SocketId = msg.SocketId
	event.Msg = string(msg.Msg)
	data, err := json.Marshal(event)
	if err != nil {
		Global.Logger.Error("Fail to marshal event:", err)
		return
	}
	for sub := Gpthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(GpManager.SocketInfo).SocketId == event.SocketId {
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
}
