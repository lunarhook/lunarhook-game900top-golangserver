package GpPacket

import (
	"container/list"
)

type EventType int

const (
	IM_C2S_JOIN = iota
	IM_S2C_JOIN
	IM_C2S_GETROOMLIST
	IM_S2C_SENDROOMLIST
	IM_C2S_JOINCREATROOM
	IM_S2C_JOINCREATROOM
	IM_S2c_BEGINQUEST
	IM_C2S_SENDQUEST
	IM_S2C_SENDQUEST
	IM_C2S_SENDANSWER
	IM_S2C_SENDANSWER
	IM_C2S_LEAVEROOM
	IM_S2C_LEAVEROOM
	IM_C2S2C_HEART
	IM_EVENT_BROADCAST_HEART
	IM_EVENT_BROADCAST_MESSAGE
	IM_EVENT_MESSAGE
)

// 用户描述信息
type IM_protocol_user struct {
	avatarUrlA string
	avatarUrlB string
}

// 用户交互协议
type IM_protocol struct {
	Type          EventType // JOIN, LEAVE, MESSAGE
	Msg           string
	SocketId      uint32
	SocketIdother uint32
	Users         IM_protocol_user
	Timestamp     int // Unix timestamp (secs)
	Roomid        uint32
}
type IM_rec struct {
	Type     EventType // JOIN, LEAVE, MESSAGE
	SocketId uint32
	Msg      string
}
type IM_ret struct {
	Type     EventType // JOIN, LEAVE, MESSAGE
	SocketId uint32
	Msg      string
}

const archiveSize = 100

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event IM_rec) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
/*
func GetEvents(lastReceived int) []IM_rec {
	events := make([]IM_rec, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(IM_rec)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
*/
