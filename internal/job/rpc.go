package job

import (
	"context"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/zlei1/goim/define"
	"github.com/zlei1/goim/proto"
	"strings"
)

var (
	RpcClientList map[string]client.XClient
)

func (j *Job) InitRpcClient() {
	d, _ := etcdClient.NewEtcdV3Discovery(j.c.Etcd.BasePath, j.c.Etcd.ServerPathComet, []string{j.c.Etcd.Host}, false, nil)
	RpcClientList = make(map[string]client.XClient, len(d.GetServices()))

	for _, cometConf := range d.GetServices() {
		cometConf.Value = strings.Replace(cometConf.Value, "=&tps=0", "", 1)
		serverId := cometConf.Value

		d, _ := client.NewPeer2PeerDiscovery(cometConf.Key, "")
		RpcClientList[serverId] = client.NewXClient(j.c.Etcd.ServerPathComet, client.Failover, client.RoundRobin, d, client.DefaultOption)
	}
}

// Push 推送单条消息
func (j *Job) Push(serverId string, userId int64, msg []byte) {
	req := &proto.PushMsgReq{
		UserId: userId,
		Msg: proto.Msg{
			Ver:       1,
			Operation: define.OP_SINGLE_SEND,
			Body:      msg,
		},
	}
	reply := &proto.SuccessReply{}

	err := RpcClientList[serverId].Call(context.Background(), "PushMsg", req, reply)
	if err != nil {
		log.Infof("job push msg err %v", err)
	}
}

// BroadcastRoom 广播到房间
func (j *Job) BroadcastRoom(roomId int64, msg []byte) {
	req := &proto.PushRoomMsgReq{
		RoomId: roomId,
		Msg: proto.Msg{
			Ver:       1,
			Operation: define.OP_ROOM_SEND,
			Body:      msg,
		},
	}
	reply := &proto.SuccessReply{}
	for _, rpc := range RpcClientList {
		err := rpc.Call(context.Background(), "PushRoomMsg", req, reply)
		if err != nil {
			log.Infof("job push room msg err %v", err)
		}
	}
}
