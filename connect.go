package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"github.com/davyxu/golog"
	"time"
	"os"
	"os/signal"
	"syscall"
)

/**
 *@author LanguageY++2013
 *2019/2/20 5:33 PM
 **/
var (
	pool *redis.Pool
	ilog = golog.New("redis_client")
)


func InitConnect(ip string, port string, password string) {
	if pool != nil {
		return
	}

	addr := ip + ":" + port

	pool =  &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			if password != "" {
				do := redis.DialPassword(password)
				c, err := redis.Dial("tcp", addr, do)
				if err != nil {
					return nil, err
				}
				return c, err
			}else{
				c, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				return c, err
			}
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	go closeConnection()
}

func closeConnection() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
		os.Exit(0)
	}()
}


func GetConn() redis.Conn{
	return pool.Get()
}