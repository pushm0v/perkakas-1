package distlock

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

type DistLock struct {
	Pool   *redis.Pool
	Option Option
}

type Option struct {
	RedisLockerTries      int
	RedisLockerRetryDelay time.Duration
	RedisLockerExpiry     time.Duration
}

const (
	RedisLockerTries      = 3
	RedisLockerRetryDelay = 2 * time.Second
	RedisLockerExpiry     = 5 * time.Second
)

var pools []redsync.Pool
var reds *redsync.Redsync

func New(pool *redis.Pool, option ...Option) (distLock *DistLock) {
	o := makeOption(option)
	return &DistLock{
		Pool:   pool,
		Option: o,
	}
}

func makeOption(option []Option) (o Option) {
	o.RedisLockerTries = RedisLockerTries
	o.RedisLockerRetryDelay = RedisLockerRetryDelay
	o.RedisLockerExpiry = RedisLockerExpiry

	for _, opt := range option {
		if opt.RedisLockerTries != 0 {
			o.RedisLockerTries = opt.RedisLockerTries
		}
		if opt.RedisLockerRetryDelay != 0 {
			o.RedisLockerRetryDelay = opt.RedisLockerRetryDelay
		}
		if opt.RedisLockerExpiry != 0 {
			o.RedisLockerExpiry = opt.RedisLockerExpiry
		}
	}
	return o
}

func (d *DistLock) SetCacheWithDistLock(key string, ttl interface{}, value interface{}) (err error) {

	pools = []redsync.Pool{d.Pool}
	reds = redsync.New(pools)

	lockKey := fmt.Sprintf("lock-key-%s", key)

	mtx := reds.NewMutex(lockKey,
		redsync.SetTries(d.Option.RedisLockerTries),
		redsync.SetRetryDelay(d.Option.RedisLockerRetryDelay),
		redsync.SetExpiry(d.Option.RedisLockerExpiry),
	)

	for {
		errLock := mtx.Lock()
		if errLock == nil {
			break
		}
	}
	defer mtx.Unlock()

	_, err = redis.Bytes(d.Pool.Get().Do("GET", key))

	if err == redis.ErrNil {
		_, err := d.Pool.Get().Do("SETEX", key, ttl, value)
		if err != nil {
			return err
		}
	}

	return nil
}
