package comet

import (
	"github.com/gorilla/websocket"
	"github.com/zlei1/goim/proto"
)

type Channel struct {
	Room *Room
	Next *Channel
	Prev *Channel
	broadcast chan *proto.Msg
	UserId int64
	conn *websocket.Conn
}

func NewChannel(size int) (c *Channel) {
	c = new(Channel)
	c.broadcast = make(chan *proto.Msg, size)
	c.Next = nil
	c.Prev = nil
	return
}

func (ch *Channel) Push(msg *proto.Msg) {
	select {
	case ch.broadcast <- msg:
	default:
	}
	return
}

