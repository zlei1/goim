package logic

import (
	"github.com/zlei1/goim/internal/logic/conf"
	"github.com/zlei1/goim/internal/logic/dao"
)

type Logic struct {
	C *conf.Config
	Dao *dao.Dao
}

func New(c *conf.Config) (l *Logic) {
	l = &Logic{
		c,
		dao.New(c),
	}
	return
}