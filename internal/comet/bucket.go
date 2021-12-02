package comet

import (
	"sync"
	"sync/atomic"

	"github.com/zlei1/goim/internal/comet/conf"
	"github.com/zlei1/goim/proto"
)

type Bucket struct {
	bc conf.Bucket
	cLock sync.RWMutex
	chs map[int64]*Channel
	rooms map[int64]*Room
	routines []chan *proto.PushRoomMsgReq
	routinesNum uint64
}

func NewBucket(bc conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.bc = bc
	b.chs = make(map[int64]*Channel, bc.ChannelSize)
	b.rooms = make(map[int64]*Room, bc.RoomSize)
	b.routines = make([]chan *proto.PushRoomMsgReq, bc.RoutineAmount)
	for i := uint64(0); i < bc.RoutineAmount; i++ {
		c := make(chan *proto.PushRoomMsgReq, bc.RoutineSize)
		b.routines[i] = c
		go b.PushRoom(c)
	}
	return
}

func (b *Bucket) Channel(UserId int64) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[UserId]
	b.cLock.RLock()
	return
}

func (b *Bucket) PushRoom(c chan *proto.PushRoomMsgReq) {
	for {
		var (
			arg  *proto.PushRoomMsgReq
			room *Room
		)
		arg = <-c
		if room = b.GetRoom(arg.RoomId); room != nil {
			room.Push(&arg.Msg)
		}
	}
}

func (b *Bucket) GetRoom(roomId int64) (room *Room) {
	b.cLock.RLock()
	room, _ = b.rooms[roomId]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) PutChannel(userId int64, roomId int64, ch *Channel) (err error) {
	var (
		room *Room
		ok bool
	)
	b.cLock.Lock()
	if roomId != NoRoom {
		if room, ok = b.rooms[roomId]; !ok {
			room = NewRoom(roomId)
			b.rooms[roomId] = room
		}
		ch.Room = room
	}
	ch.UserId = userId
	b.chs[userId] = ch
	b.cLock.Unlock()

	if room != nil {
		err = room.PutChannel(ch)
	}
	return
}

func (b *Bucket) DeleteChannel(ch *Channel) {
	var (
		ok bool
		room *Room
	)
	b.cLock.RLock()
	if ch, ok = b.chs[ch.UserId]; ok {
		room = b.chs[ch.UserId].Room
		delete(b.chs, ch.UserId)
	}
	b.cLock.RUnlock()
	if room != nil && room.DeleteChannel(ch) {
		if room.drop == true {
			delete(b.rooms, room.Id)
		}
	}
}

func (b *Bucket) BroadcastRoom(req *proto.PushRoomMsgReq) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.bc.RoutineAmount
	b.routines[num] <- req
}