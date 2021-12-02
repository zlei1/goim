package job

import (
	"github.com/go-redis/redis"
	"github.com/zlei1/goim/define"
)

var RedisClient *redis.Client

func (j *Job) InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     j.c.Redis.Addr,
		Password: j.c.Redis.Password,        // no password set
		DB:       j.c.Redis.DB, // use default DB
	})
	if _, err := RedisClient.Ping().Result(); err != nil {
		panic(err)
	}

	go func() {
		redisSub := RedisClient.Subscribe(define.REDIS_CHANNEL)
		ch := redisSub.Channel()
		for {
			msg, ok := <-ch
			if !ok {
				break
			}

			redisPush(j, msg.Payload)
		}
	}()
}