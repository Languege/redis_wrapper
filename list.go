package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/
func LPush(key string, value []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", key, value)

	return err
}

func RPop(key string) ([]byte, error) {

	conn := pool.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("RPOP", key))
}