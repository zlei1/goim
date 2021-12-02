package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	confPath string
	Conf *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "comet-example.toml", "default config path")
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
			RPCAddress: "localhost:3020",
			WriteWait: 10 * time.Second,
			PongWait: 60 * time.Second,
			PingPeriod: 54 * time.Second,
			MaxMessageSize: 512,
			BroadcastSize: 512,
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
		Bucket: Bucket{
			Num: 8,
			ChannelSize: 1024,
			RoomSize: 1024,
			RoutineAmount: 32,
			RoutineSize: 20,
		},
		Websocket: Websocket{
			Bind: "0.0.0.0:8911",
		},
	}
}

type Config struct {
	Base Base
	Bucket Bucket
	Etcd Etcd
	Websocket Websocket
}

type Base struct {
	ServerId string
	RPCAddress string
	WriteWait time.Duration
	PongWait time.Duration
	PingPeriod time.Duration
	MaxMessageSize int64
	BroadcastSize int
	ReadBufferSize int
	WriteBufferSize int
}

type Bucket struct {
	Num int
	ChannelSize int
	RoomSize int
	RoutineAmount uint64
	RoutineSize int
}

type Etcd struct {
	Host string
	BasePath string
	ServerId string
	ServerPath string
}

type Websocket struct {
	Bind string
}