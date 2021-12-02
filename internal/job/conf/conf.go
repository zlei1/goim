package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

var (
	confPath string
	Conf *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "job-example.toml", "default config path")
}

func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	log.Infof("job config: %+v", Conf)
	return
}

func Default() *Config {
	return &Config{
		Base: Base{
			RPCAddr: "localhost:3001",
			PushChan: 2,
			PushSize: 50,
		},
	}
}

type Config struct {
	Base Base
	Etcd Etcd
	Redis Redis
}

type Base struct {
	RPCAddr string
	PushChan int
	PushSize int
}

type Etcd struct {
	Host string
	BasePath string
	ServerPath string
}

type Redis struct {
	Addr string
	Password string
	DB int
}