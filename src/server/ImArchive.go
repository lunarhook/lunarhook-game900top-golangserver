// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package gameserver

import (
	"container/list"
)

type EventType int

const (
	IM_EVENT_JOIN = iota
	IM_EVENT_BROADCAST_JOIN
	IM_EVENT_LEAVE
	IM_EVENT_BROADCAST_LEAVE
	IM_EVENT_MESSAGE
	IM_EVENT_BROADCAST_MESSAGE
	IM_EVENT_HEART
	IM_EVENT_BROADCAST_HEART
)


//用户描述信息
type IM_protocol_user struct {
	To 			uint32
	From		uint32
	SessKey		uint32
	ChatChanId	uint32
}

//用户交互协议
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
