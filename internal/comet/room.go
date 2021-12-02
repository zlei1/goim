package comet

import (
	"errors"
	"github.com/zlei1/goim/proto"
	"sync"
)

const NoRoom = 0

type Room struct {
	Id int64
	OnlineCount int
	rLock sync.RWMutex
	drop bool
	next *Channel
}

func NewRoom(roomId int64) *Room {
	room := new(Room)
	room.Id = roomId
	room.drop = false
	room.next = nil
	room.OnlineCount = 0
	return room
}

func (r *Room) PutChannel(ch *Channel) (err error) {
	r.rLock.Lock()
	if r.drop {
		err = errors.New("room drop")
	} else {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.OnlineCount++
	}
	r.rLock.Unlock()
	return
}

func (r *Room) DeleteChannel(ch *Channel) bool {
	r.rLock.Lock()
	if ch.Prev == nil && ch.Next == nil {
		r.rLock.Unlock()
		return false
	}
	if ch.Next != nil {
		// if not footer
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.OnlineCount--
	r.drop = r.OnlineCount == 0
	r.rLock.Unlock()
	return r.drop
}

func (r *Room) Push(msg *proto.Msg) {
	r.rLock.Lock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Push(msg)
	}
	r.rLock.Unlock()
	return
}