package Gamecontrollers

import (
	GpPacket "github.com/lunarhook/lunarhook-game900top-golangserver/src/server"
	Global "github.com/lunarhook/lunarhook-game900top-golangserver/src/server/Global"
)

func (this *WebSocketListController) chatroom() {
	for {
		select {
		case JoinSocket := <-this.SocketChan:
			if !this.IsExistSocketById(JoinSocket.SocketId) {
				this.ActiveSocketList.PushBack(JoinSocket) // Add user to the end of list.
				// Publish a JOIN event.
				this.MsgList <- this.NewMsg(GpPacket.IM_S2C_JOIN, JoinSocket.User, JoinSocket.SocketId, "")

				this.MsgList <- this.NewMsg(GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("New socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			} else {
				this.MsgList <- this.NewMsg(GpPacket.IM_EVENT_MESSAGE, JoinSocket.User, JoinSocket.SocketId, "welcome")
				Global.Logger.Info("Old socket:", JoinSocket.SocketId, ";WebSocket:", JoinSocket.Conn != nil)
			}
		case SocketMessage := <-this.MsgList:
			//如果是心跳，单发
			switch SocketMessage.Type {
			case
				GpPacket.IM_C2S2C_HEART: // 心跳
				this.HeartWebSocket(SocketMessage)
				break
			case
				GpPacket.IM_S2C_LEAVE: // 关闭房间完成比赛，退出游戏，展示结果
				break
			case
				GpPacket.IM_S2C_JOIN, //创建房间，显示房间列表
				GpPacket.IM_EVENT_MESSAGE:
				this.HeartWebSocket(SocketMessage)
				this.Game(SocketMessage)
				break
			case
				GpPacket.IM_S2C_CREATROOM: //创建房间等待对手加入
				break
			case
				GpPacket.IM_S2C_ROOMLIST: // 刷新房间列表
				this.broadcastWebSocket(SocketMessage)
				break
			case
				GpPacket.IM_EVENT_BROADCAST_HEART,

				GpPacket.IM_EVENT_BROADCAST_MESSAGE:
				this.broadcastWebSocket(SocketMessage)
				break
			}
			GpPacket.NewArchive(SocketMessage)
			if SocketMessage.Type == GpPacket.IM_EVENT_MESSAGE {
				Global.Logger.Info("Message from", SocketMessage.Users.From, ";Msg:", SocketMessage.Msg)
			}

		case LeaveSocket := <-this.UnSocketChan:
			for sub := this.ActiveSocketList.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(SocketInfo).SocketId == LeaveSocket.SocketId {
					this.ActiveSocketList.Remove(sub)
					// Clone connection.
					ws := sub.Value.(SocketInfo).Conn
					if ws != nil {
						ws.Close()
						Global.Logger.Error("WebSocket closed:", LeaveSocket)
					}

					this.MsgList <- this.NewMsg(GpPacket.IM_S2C_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "") // Publish a LEAVE event.
					//this.MsgList <- this.NewMsg(GpPacket.IM_EVENT_BROADCAST_LEAVE, sub.Value.(SocketInfo).User, LeaveSocket.SocketId, "")
					break
				}
			}
		}
	}
}
