package redis_wrapper

import "github.com/gomodule/redigo/redis"

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/
func HSet(key,field string, value []byte) (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("HSET", key, field, value))
}

func HGet(key,field string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("HGET", key, field))
}

//data = 1存在，data = 0不存在
func HExist(key, field string)(int64, error) {
	conn := pool.Get()
	defer conn.Close()

	data, err := redis.Int64(conn.Do("HEXISTS", key, field))

	return data, err
}

func HDel(key, field string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", key, field)

	return err
}