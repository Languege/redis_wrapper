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


func InitConnect(ip, port, password string, options... interface{}) {
	if wrapper != nil {
		return
	}

	wrapper = NewRedisWrapper(ip, port, password, 100, time.Duration(240) * time.Second, 200)
	if len(options) > 0 {
		prefix, ok := options[0].(string)
		if ok {
			wrapper.Prefix = prefix
		}
	}
}


func GetConn() redis.Conn{
	return wrapper.Get()
}