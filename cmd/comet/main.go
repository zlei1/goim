package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/internal/comet"
	"github.com/zlei1/goim/internal/comet/conf"
)

const (
	ver = "1.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Infof("comet server [version: %s] start", ver)

	s := comet.New(conf.Conf)
	s.InitLogicRpcClient()
	s.InitRpcServer()
	s.InitWebsocket()
}