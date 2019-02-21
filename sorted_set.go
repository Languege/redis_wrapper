package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:32 PM
 **/

 func ZAdd(key string, score float64, value []byte)  error {
 	conn := pool.Get()
 	defer conn.Close()

 	var err error
 	_, err = conn.Do("ZADD", key, score, value)
 	return err
 }

 func ZCard(key string) (size int64, err error) {
	 conn := pool.Get()
	 defer conn.Close()

	 return redis.Int64(conn.Do("ZCard", key))
 }


 //根据score获取数据
func ZRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	conn := pool.Get()
	defer conn.Close()

	if withScores {
		if count > 0 {
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
		}
	}else{
		if count > 0 {
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max))
		}
	}


	return
}
