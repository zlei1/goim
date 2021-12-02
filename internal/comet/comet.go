package comet

import (
	"github.com/zhenjl/cityhash"
	"github.com/zlei1/goim/internal/comet/conf"
)

type Comet struct {
	C *conf.Config
	Buckets []*Bucket
	BucketsLen uint32
	operator Operator
}

func New(c *conf.Config) (s *Comet) {
	buckets := make([]*Bucket, c.Bucket.Num)
	for i := 0; i < c.Bucket.Num; i++ {
		buckets[i] = NewBucket(c.Bucket)
	}

	s = &Comet{
		c,
		buckets,
		uint32(len(buckets)),
		new(DefaultOperator),
	}
	return
}

func (cmt *Comet) assignBucket(key string) *Bucket {
	index := cityhash.CityHash32([]byte(key), uint32(len(key))) % cmt.BucketsLen
	return cmt.Buckets[index]
}