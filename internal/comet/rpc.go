package comet

import (
	"context"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/zlei1/goim/define"
	"github.com/zlei1/goim/proto"
	"strconv"
	"strings"
	"time"
)

type Rpc struct {
	cmt *Comet
}

func (cmt *Comet) InitRpcServer() {
	address := strings.Split(cmt.C.Base.RPCAddress, ",")
	for _, addr := range address {
		go createServer(addr, cmt)
	}
}

func createServer(addr string, cmt *Comet) {
	s := server.NewServer()
	addRegistryPlugin(addr, s, cmt)
	err := s.RegisterName(cmt.C.Etcd.ServerPath, &Rpc{cmt: cmt}, cmt.C.Etcd.ServerId)
	if err != nil {
		panic(err)
	}
	err = s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func addRegistryPlugin(addr string, s *server.Server, cmt *Comet) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		EtcdServers:    []string{cmt.C.Etcd.Host},
		BasePath:       cmt.C.Etcd.BasePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	s.Plugins.Add(r)
}

func (rpc *Rpc) PushMsg(ctx context.Context, req *proto.PushMsgReq, reply *proto.SuccessReply) (err error) {
	var (
		bucket  *Bucket
		channel *Channel
	)
	bucket = rpc.cmt.assignBucket(strconv.FormatInt(req.UserId, 10))
	if channel = bucket.Channel(req.UserId); channel != nil {
		channel.Push(&req.Msg)
		return
	}

	reply.Code = define.SUCCESS_REPLY
	reply.Msg = define.SUCCESS_REPLY_MSG
	return nil
}

func (rpc *Rpc) PushRoomMsg(ctx context.Context, req *proto.PushRoomMsgReq, reply *proto.SuccessReply) (err error) {
	reply.Code = define.SUCCESS_REPLY
	reply.Msg = define.SUCCESS_REPLY_MSG
	for _, bucket := range rpc.cmt.Buckets {
		bucket.BroadcastRoom(req)
	}
	return nil
}