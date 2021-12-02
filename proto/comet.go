package proto

type Msg struct {
	Ver int `json:"ver"`
	Operation int32 `json:"op"`
	Body []byte `json:"body"`
}

type PushMsgReq struct {
	UserId int64
	Msg Msg
}

type PushRoomMsgReq struct {
	RoomId int64
	Msg Msg
}