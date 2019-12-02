package distlock

import (
	"github.com/gomodule/redigo/redis"
	"github.com/im7mortal/kmutex"
)

type DistLock struct {
	Pool   *redis.Pool
	Kmutex *kmutex.Kmutex
}

func New(pool *redis.Pool) (distLock *DistLock) {
	kMutex := kmutex.New()

	return &DistLock{
		Pool:   pool,
		Kmutex: kMutex,
	}
}

func (d *DistLock) SetCacheWithDistLock(key string, ttl interface{}, value interface{}) (err error) {

	d.Kmutex.Lock(key)
	defer d.Kmutex.Unlock(key)

	_, err = redis.Bytes(d.Pool.Get().Do("GET", key))
	if err == redis.ErrNil {
		_, err := d.Pool.Get().Do("SETEX", key, ttl, value)
		if err != nil {
			return err
		}
	}

	return nil
}
