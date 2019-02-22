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
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES"))
		}
	}else{
		if count > 0 {
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max, "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", key, min, max))
		}
	}


	return
}

func ZRevRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	conn := pool.Get()
	defer conn.Close()

	if withScores {
		if count > 0 {
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, min, max, "WITHSCORES", "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, min, max, "WITHSCORES"))
		}
	}else{
		if count > 0 {
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, min, max, "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", key, min, max))
		}
	}


	return
}

 func ZRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	 conn := pool.Get()
	 defer conn.Close()

	 if withScores {
	 	_, err = conn.Do("ZRANGE", key, start, stop, "WITHSCORES")
	 }else{
		 _, err = conn.Do("ZRANGE", key, start, stop)
	 }

 	return
 }

func ZRevRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()

	if withScores {
		_, err = conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES")
	}else{
		_, err = conn.Do("ZREVRANGE", key, start, stop)
	}

	return
}

 func ZIncreBy(key string, increment float64, member interface{})(err error) {
	 conn := pool.Get()
	 defer conn.Close()

	 _, err = conn.Do("ZINCRBY", key, increment, member)
	 return
 }

 //移除一个元素
 func ZRem(key string, member interface{})(err error) {
	 conn := pool.Get()
	 defer conn.Close()

	 _, err = conn.Do("ZREM", key, member)
	 return
 }
