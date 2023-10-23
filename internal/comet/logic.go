package comet

import (
	"context"
	etcdClient "github.com/rpcxio/rpcx-etcd/client"
	log "github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/zlei1/goim/proto"
)

var logicRpcClient client.XClient

func (cmt *Comet) InitLogicRpcClient() {
	d, _ := etcdClient.NewEtcdV3Discovery(
		cmt.C.Etcd.BasePath,
		cmt.C.Etcd.ServerPathLogic,
		[]string{cmt.C.Etcd.Host},
		false,
		nil,
	)
	logicRpcClient = client.NewXClient(cmt.C.Etcd.ServerPathLogic, client.Failtry, client.RandomSelect, d, client.DefaultOption)
}

func LogicConnect(req *proto.ConnectReq) (err error) {
	res := &proto.ConnectRes{}
	err = logicRpcClient.Call(context.Background(), "Connect", req, res)
	if err != nil {
		log.Errorf("comet logicRpcClient call connect err: %s", err.Error())
	}
	return
}

func LogicDisconnect(req *proto.DisconnectReq) (err error) {
	res := &proto.DisconnectRes{}
	err = logicRpcClient.Call(context.Background(), "Disconnect", req, res)
	if err != nil {
		log.Errorf("comet logicRpcClient call disconnect err: %s", err.Error())
	}
	return
}
