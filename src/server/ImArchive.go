package GpPacket

import (
	"container/list"
)

type EventType int

const (
	IM_S2C_JOIN = iota
	IM_S2C_ROOMLIST
	IM_S2C_CREATROOM
	IM_EVENT_MESSAGE
	IM_EVENT_BROADCAST_MESSAGE
	IM_C2S2C_HEART
	IM_S2C_LEAVE
	IM_EVENT_BROADCAST_HEART
)

// 用户描述信息
type IM_protocol_user struct {
	To         uint32
	From       uint32
	SessKey    uint32
	ChatChanId uint32
}

// 用户交互协议
type IM_protocol struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	Msg       string
	SocketId  uint32
	Users     IM_protocol_user
	Timestamp int // Unix timestamp (secs)
}

const archiveSize = 100

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event IM_protocol) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(lastReceived int) []IM_protocol {
	events := make([]IM_protocol, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(IM_protocol)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
