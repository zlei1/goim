package comet

import "github.com/zlei1/goim/proto"

type Operator interface {
	Connect(req *proto.ConnectReq) error
	Disconnect(req *proto.DisconnectReq) error
}

type DefaultOperator struct {
}

func (o *DefaultOperator) Connect(req *proto.ConnectReq) (err error) {
	err = LogicConnect(req)
	if err != nil {
		return err
	}
	return nil
}

func (o *DefaultOperator) Disconnect(req *proto.DisconnectReq) (err error) {
	err = LogicDisconnect(req)
	if err != nil {
		return err
	}
	return nil
}