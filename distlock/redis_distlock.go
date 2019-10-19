package distlock

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
)

const (
	RedisLockerTries      = 3
	RedisLockerRetryDelay = 2 * time.Second
	RedisLockerExpiry     = 5 * time.Second
)

var pools []redsync.Pool
var reds *redsync.Redsync

func SetCacheWithDistLock(pool *redis.Pool, key string, ttl interface{}, value interface{}) (err error) {

	pools = []redsync.Pool{pool}
	reds = redsync.New(pools)

	lockKey := fmt.Sprintf("lock-key-%s", key)

	mtx := reds.NewMutex(lockKey,
		redsync.SetTries(RedisLockerTries),
		redsync.SetRetryDelay(RedisLockerRetryDelay),
		redsync.SetExpiry(RedisLockerExpiry),
	)

	for {
		errLock := mtx.Lock()
		if errLock == nil {
			break
		}
	}
	defer mtx.Unlock()

	_, err = redis.Bytes(pool.Get().Do("GET", key))

	if err == redis.ErrNil {
		_, err := pool.Get().Do("SETEX", key, ttl, value)
		if err != nil {
			return err
		}
	}

	return nil
}
