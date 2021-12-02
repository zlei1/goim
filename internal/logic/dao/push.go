package dao

import (
	"encoding/json"
	"github.com/zlei1/goim/define"
	"github.com/zlei1/goim/proto"
)

func (d *Dao) PushMsg(serverId string, userId int64, msg []byte) (err error) {
	var redisMsg = &proto.JobMessage{
		Op: define.OP_SINGLE_SEND,
		ServerId: serverId,
		UserId: userId,
		Msg: msg,
	}
	redisMsgByte, err := json.Marshal(redisMsg)
	err = d.Redis.Publish(define.REDIS_CHANNEL, redisMsgByte).Err()
	return
}

func (d *Dao) BroadcastRoomMsg(roomId int64, msg []byte) (err error) {
	var redisMsg = &proto.JobMessage{
		Op: define.OP_ROOM_SEND,
		RoomId: roomId,
		Msg: msg,
	}
	redisMsgByte, err := json.Marshal(redisMsg)
	err = d.Redis.Publish(define.REDIS_CHANNEL, redisMsgByte).Err()
	return
}

func (d *Dao) BroadcastMsg() (err error) {
	return
}
