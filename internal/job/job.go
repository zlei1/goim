package job

import (
	"github.com/zlei1/goim/internal/job/conf"
)

type Job struct {
	c *conf.Config
}

func New(c *conf.Config) (j *Job) {
	j = &Job{
		c,
	}
	return
}