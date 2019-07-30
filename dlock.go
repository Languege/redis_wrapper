package redis_wrapper

import (
	"time"
	"github.com/gomodule/redigo/redis"
)

/**
 *@author LanguageY++2013
 *2019/2/22 6:36 PM
 **/
//分布式锁

func TryLock(key string, seconds int)(uniqueID int64, err error) {
	conn := pool.Get()
	defer conn.Close()

	uniqueID = time.Now().UnixNano()
	_, err = redis.String(conn.Do("SET", key, uniqueID, "EX", seconds, "NX"))
	return
}

func Release(key string, uniqueID int64)(err error) {
	conn := pool.Get()
	defer conn.Close()


	script := redis.NewScript(1, `
	 if redis.call("get",KEYS[1]) == ARGV[1] then
          return redis.call("del",KEYS[1])
      else
          return 0
      end
	`)

	_, err = script.Do(conn, key, uniqueID)

	return
}

