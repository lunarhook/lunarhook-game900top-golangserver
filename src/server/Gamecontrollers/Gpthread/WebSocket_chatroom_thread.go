package Gpthread

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers"
	"github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Gamecontrollers/Gphandle"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func Chatroom(Gthis *Gamecontrollers.WebSocketListController) {
	for {
		select {
		case JoinSocket := <-Gthis.SocketChan:
			if !Gthis.IsExistSocketById(JoinSocket.SocketId) {
				Gthis.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_S2C_JOIN, JoinSocket.User, JoinSocket.SocketId, "")

				Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("New socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			} else {
				Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("Old socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			}
		case SocketMessage := <-Gthis.MsgList:
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_C2S2C_HEART: // 心跳
				Gphandle.HeartWebSocket(SocketMessage)
				break
			case
				GpPacket.IM_S2C_LEAVE: // 关闭房间完成比赛，退出游戏，展示结果
				break
			case
				GpPacket.IM_S2C_JOIN, //创建房间，显示房间列表
				GpPacket.IM_EVENT_MESSAGE:
				Gphandle.HeartWebSocket(SocketMessage, Gamecontrollers.GlobaWebSocketListManager)
				Gthis.Game(SocketMessage)
				break
			case
				GpPacket.IM_S2C_CREATROOM: //创建房间等待对手加入
				break
			case
				GpPacket.IM_S2C_ROOMLIST: // 刷新房间列表
				Gphandle.BroadcastWebSocket(SocketMessage, Gamecontrollers.GlobaWebSocketListManager)
				break
			case
				GpPacket.IM_EVENT_BROADCAST_HEART,

				GpPacket.IM_EVENT_BROADCAST_MESSAGE:
				Gphandle.BroadcastWebSocket(SocketMessage, Gamecontrollers.GlobaWebSocketListManager)
				break
			}
			GpPacket.NewArchive(SocketMessage)
			if SocketMessage.Type == GpPacket.IM_EVENT_MESSAGE {
				Global.Logger.Info("Message from", SocketMessage.Users.From, ";Msg:", SocketMessage.Msg)
			}

		case LeaveSocket := <-Gthis.UnSocketChan:
			for sub := Gthis.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Gamecontrollers.SocketInfo).SocketId == LeaveSocket.SocketId {
					Gthis.ActiveSocketList.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Gamecontrollers.SocketInfo).Conn
					if ws != nil {
						ws.Close()
						Global.Logger.Error("WebSocket closed:", LeaveSocket)
					}

					Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_S2C_LEAVE, sub.Value.(Gamecontrollers.SocketInfo).User, LeaveSocket.SocketId, "") // Publish a LEAVE event.
					//Gthis.MsgList <- Gthis.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
					break
				}
			}
		}
	}
}