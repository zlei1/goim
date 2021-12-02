package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/internal/logic"
	"github.com/zlei1/goim/internal/logic/conf"
	"github.com/zlei1/goim/internal/logic/http"
	"github.com/zlei1/goim/internal/logic/rpc"
	"os"
	"os/signal"
	"syscall"
)

const (
	ver = "1.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Infof("logic server [version: %s] start", ver)

	l := logic.New(conf.Conf)
	http.InitHTTP(l)
	rpc.InitLogicRpcServer(l)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("Shutting down server...")
}