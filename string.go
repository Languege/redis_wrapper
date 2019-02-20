package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
	"strconv"
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

	commandName := fmt.Sprintf("SET %s %s", key, string(value))
	if ex > 0 {
		commandName += " EX " + strconv.Itoa(ex)
	}

	if px > 0 {
		commandName += " EX " + strconv.Itoa(px)
	}

	if nx {
		commandName += " NX"
	}else if xx {
		commandName += " XX"
	}

	_, err = conn.Do("SET", key, value)

	return err
}

func Incr(key string) (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("INCR", key))
}