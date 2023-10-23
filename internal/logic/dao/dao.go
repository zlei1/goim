package dao

import (
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/zlei1/goim/define"
	"github.com/zlei1/goim/internal/logic/conf"
)

type Dao struct {
	c     *conf.Config
	Redis *redis.Client
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		c:     c,
		Redis: newRedis(c.Redis),
	}
	return d
}

func newRedis(c conf.Redis) *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	if _, err := RedisClient.Ping().Result(); err != nil {
		log.Panic(err)
	}

	return RedisClient
}

func (d *Dao) GetServerByUserId(userId int64) (serverId string, err error) {
	serverId, err = d.Redis.Get(userKey(userId)).Result()
	if err != nil {
		return
	}
	return
}

func (d *Dao) SetServerWithUserId(userId int64, serverId string) (err error) {
	err = d.Redis.Set(userKey(userId), serverId, 0).Err()
	return
}

func userKey(userId int64) string {
	return fmt.Sprintf(define.PREFIX_USERID, userId)
}
