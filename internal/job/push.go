package job

import (
	"encoding/json"
	"github.com/zlei1/goim/define"
	"github.com/zlei1/goim/proto"
	"math/rand"
)

type PushParams struct {
	ServiceId string
	UserId int64
	RoomId int64
	Msg []byte
}

var pushChannels []chan *PushParams

func (j *Job) InitPush() {
	pushChannels = make([]chan *PushParams, j.c.Base.PushChan)
	for i := 0; i < len(pushChannels); i++ {
		pushChannels[i] = make(chan *PushParams, j.c.Base.PushSize)
		go processSinglePush(j, pushChannels[i])
	}
}

func processSinglePush(j *Job, ch chan *PushParams) {
	var params *PushParams
	for {
		params = <-ch
		j.Push(params.ServiceId, params.UserId, params.Msg)
	}
}

func redisPush(j *Job, msg string) (err error) {
	m := &proto.JobMessage{}
	if err := json.Unmarshal([]byte(msg), m); err != nil {
		return err
	}

	switch m.Op {
	case define.OP_SINGLE_SEND:
		pushChannels[rand.Int()%j.c.Base.PushChan] <- &PushParams{
			ServiceId: m.ServerId,
			UserId: m.UserId,
			Msg: m.Msg,
		}
		break
	case define.OP_ROOM_SEND:
		j.BroadcastRoom(m.RoomId, m.Msg)
		break
	}
	return
}
