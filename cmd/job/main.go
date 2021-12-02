package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/internal/job"
	"github.com/zlei1/goim/internal/job/conf"
)

const (
	ver = "1.0"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Infof("job server [version: %s] start", ver)

	j := job.New(conf.Conf)
	j.InitRedis()
	j.InitRpcClient()
	j.InitPush()

	select {}
}