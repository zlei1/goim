package rpc

import (
	"context"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/zlei1/goim/internal/logic"
	"github.com/zlei1/goim/proto"
	"strings"
	"time"
)

type LogicRpcServer struct {
	l *logic.Logic
}

func InitLogicRpcServer(l *logic.Logic) {
	address := strings.Split(l.C.Base.RPCAddress, ",")
	for _, addr := range address {
		go createServer(addr, l)
	}
}

func createServer(addr string, l *logic.Logic) {
	s := server.NewServer()
	addRegistryPlugin(addr, s, l)
	s.RegisterName(l.C.Etcd.ServerPath, &LogicRpcServer{l: l}, l.C.Etcd.ServerId)
	s.Serve("tcp", addr)
}

func addRegistryPlugin(addr string, s *server.Server, l *logic.Logic) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + addr,
		EtcdServers:    []string{l.C.Etcd.Host},
		BasePath:       l.C.Etcd.BasePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	s.Plugins.Add(r)
}

func (rpc *LogicRpcServer) Connect(ctx context.Context, req *proto.ConnectReq, res *proto.ConnectRes) (err error) {
	res.UserId = req.UserId
	if err = rpc.l.Dao.SetServerWithUserId(req.UserId, req.ServerId); err != nil {
		return err
	}
	return nil
}

func (rpc *LogicRpcServer) Disconnect(ctx context.Context, req *proto.DisconnectReq, res *proto.DisconnectRes) (err error) {
	return nil
}
