package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:33 PM
 **/
const(
	Default_MaxIdle = 10
	Default_IdleTimeout = time.Duration(240) * time.Second
	Default_MaxActive = 20
)
var(
	wrapper	*RedisWrapper
)


func InitConnect(ip, port, password string,maxIdle, maxActive int, idleTimeout time.Duration, options... interface{}) {
	if wrapper != nil {
		return
	}

	if maxIdle == 0 {
		maxIdle = Default_MaxIdle
	}

	if idleTimeout == 0 {
		idleTimeout = Default_IdleTimeout
	}

	if maxActive == 0 {
		maxActive = Default_MaxActive
	}

	wrapper = NewRedisWrapper(ip, port, password, maxIdle, idleTimeout, maxActive)
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