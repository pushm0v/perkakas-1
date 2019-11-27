package distlock

import (
	"fmt"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/im7mortal/kmutex"
)

type DistLock struct {
	Pool *redis.Pool
}

var pools []redsync.Pool
var reds *redsync.Redsync
var kMutex *kmutex.Kmutex

func New(pool *redis.Pool) (distLock *DistLock) {
	pools = []redsync.Pool{pool}
	reds = redsync.New(pools)
	kMutex = kmutex.New()

	return &DistLock{
		Pool: pool,
	}
}

func (d *DistLock) SetCacheWithDistLock(key string, ttl interface{}, value interface{}) (err error) {

	lockKey := fmt.Sprintf("lock-key-%s", key)

	kMutex.Lock(lockKey)
	defer kMutex.Unlock(lockKey)

	_, err = redis.Bytes(d.Pool.Get().Do("GET", key))

	if err == redis.ErrNil {
		_, err := d.Pool.Get().Do("SETEX", key, ttl, value)
		if err != nil {
			return err
		}
	}

	return nil
}
