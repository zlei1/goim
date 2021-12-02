package proto

type ConnectReq struct {
	UserId int64
	RoomId int64
	ServerId string
}

type ConnectRes struct {
	UserId int64
}

type DisconnectReq struct {
	RoomId int64
	UserId int64
}

type DisconnectRes struct {
	Has bool
}