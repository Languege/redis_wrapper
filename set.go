package redis_wrapper

import "github.com/gomodule/redigo/redis"

/**
 *@author LanguageY++2013
 *2019/2/20 5:32 PM
 **/
func SAdd(key string, member interface{}) (err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SADD", key, member)
	return
}

func SRem(key string, member interface{})(err error) {
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("SREM", key, member)
	return
}

func SCard(key string)(size int, err error) {
	conn := pool.Get()
	defer conn.Close()

	size, err = redis.Int(conn.Do("SREM", key))
	return
}

func SPop(key string)(value interface{}, err error){
	conn := pool.Get()
	defer conn.Close()

	value, err = redis.Int(conn.Do("SPOP", key))
	return
}

func SMembers(key string)(values []interface{}, err error) {
	conn := pool.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("SMEMBERS", key))
	return
}

func SRandMember(key string, count int)(values []interface{}, err error){
	conn := pool.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("SRANDMEMBER", key, count))
	return
}

func SIsMember(key string, member interface{})(value bool, err error){
	conn := pool.Get()
	defer conn.Close()

	value, err = redis.Bool(conn.Do("SISMEMBER", key, member))
	return
}
