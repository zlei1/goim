package logic

func (l *Logic) Push(userId int64, msg []byte) (err error) {
	serverId, err := l.Dao.GetServerByUserId(userId)
	if err != nil {
		return
	}
	return l.Dao.PushMsg(serverId, userId, msg)
}

func (l *Logic) PushRoom(roomId int64, msg []byte) (err error) {
	return l.Dao.BroadcastRoomMsg(roomId, msg)
}