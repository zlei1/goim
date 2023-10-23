package conf

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/define"
	"time"
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "logic-example.toml", "default config path")
}

func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	log.Infof("comet config: %+v", Conf)
	return
}

func Default() *Config {
	return &Config{
		Base: Base{
			HttpAddr:         "localhost:3000",
			RPCAddress:       "localhost:3030",
			ServerId:         genServerId(),
			HTTPReadTimeout:  10 * time.Second,
			HTTPWriteTimeout: 20 * time.Second,
		},
	}
}

type Config struct {
	Base  Base
	Etcd  Etcd
	Redis Redis
}

type Base struct {
	HttpAddr         string
	RPCAddress       string
	ServerId         string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
}

type Etcd struct {
	Host            string
	BasePath        string
	ServerPathLogic string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func genServerId() string {
	return fmt.Sprintf(define.PREFIX_LOGIC, uuid.New())
}
