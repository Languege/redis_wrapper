package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/

func Get(key string) ([]byte, error) {

	conn := pool.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("GET", key))
}

func Set(key string, value []byte, ex int, px int, nx bool, xx bool) error {
	conn := pool.Get()
	defer conn.Close()

	var err error

	if ex > 0 {
		if nx {
			_, err = conn.Do("SET", key, value, "EX", ex, "NX")
		}else if xx {
			_, err = conn.Do("SET", key, value, "EX", ex, "XX")
		}else{
			_, err = conn.Do("SET", key, value, "EX", ex)
		}
	}else if px > 0 {
		if nx {
			_, err = conn.Do("SET", key, value, "PX", px, "NX")
		}else if xx {
			_, err = conn.Do("SET", key, value, "PX", px, "XX")
		}else{
			_, err = conn.Do("SET", key, value, "PX", px)
		}
	}else{
		if nx {
			_, err = conn.Do("SET", key, value, "NX")
		}else if xx {
			_, err = conn.Do("SET", key, value, "XX")
		}else{
			_, err = conn.Do("SET", key, value)
		}
	}

	return err
}

func Incr(key string) (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("INCR", key))
}