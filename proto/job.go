package proto

type JobMessage struct {
	Op int32 `json:"op"`
	ServerId string `json:"server_id,omitempty"`
	RoomId int64 `json:"room_id,omitempty"`
	UserId int64 `json:"user_id,omitempty"`
	Msg []byte `json:"msg"`
}