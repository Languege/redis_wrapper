package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:33 PM
 **/
var(
	wrapper	*RedisWrapper
)


func InitConnect(ip string, port string, password string) {
	if wrapper != nil {
		return
	}

	wrapper = NewRedisWrapper(ip, port, password, 100, time.Duration(240) * time.Second, 200)
}


func GetConn() redis.Conn{
	return wrapper.Get()
}