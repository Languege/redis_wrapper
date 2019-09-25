package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"os"
	"os/signal"
	"syscall"
)

/**
 *@author LanguageY++2013
 *2019/8/9 9:12 PM
 **/
type RedisWrapper struct {
	redis.Pool
	Prefix 		string	//key前缀
}

func NewRedisWrapper(ip string, port string, password string, maxIdle int, idleTimeout time.Duration, maxActive int) *RedisWrapper{
	addr := ip + ":" + port

	w  :=  &RedisWrapper{
		Pool:redis.Pool{
			MaxIdle:     maxIdle,
			IdleTimeout: idleTimeout,

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
			MaxActive:maxActive,
		},
	}

	go w.closeConnection()

	return w
}

func(self *RedisWrapper) buildKey(key string) string {
	if self.Prefix != "" {
		key = self.Prefix + string(os.PathListSeparator) + key
	}

	return key
}

func(self *RedisWrapper) closeConnection() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		self.Close()
		os.Exit(0)
	}()
}

func(self *RedisWrapper) GetConn() redis.Conn{
	return self.Get()
}

func(self *RedisWrapper) TryLock(key string, seconds int)(uniqueID int64, err error) {
	conn := self.Get()
	defer conn.Close()

	uniqueID = time.Now().UnixNano()
	_, err = redis.String(conn.Do("SET", self.buildKey(key), uniqueID, "EX", seconds, "NX"))
	return
}

func(self *RedisWrapper) Release(key string, uniqueID int64)(err error) {
	conn := self.Get()
	defer conn.Close()


	script := redis.NewScript(1, `
	 if redis.call("get",KEYS[1]) == ARGV[1] then
          return redis.call("del",KEYS[1])
      else
          return 0
      end
	`)

	_, err = script.Do(conn, self.buildKey(key), uniqueID)

	return
}

func(self *RedisWrapper) HSet(key,field string, value []byte) (int64, error) {
	conn := self.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("HSET", self.buildKey(key), field, value))
}

func(self *RedisWrapper) HMSet(key string, kv map[string]string) (string, error) {
	conn := self.Get()
	defer conn.Close()

	params := make([]interface{}, 0, 2 * len(kv)+1)
	params = append(params, self.buildKey(key))
	for k, v := range kv  {
		params = append(params, k, v)
	}

	return redis.String(conn.Do("HMSET",  params...))
}

func(self *RedisWrapper) HGet(key,field string) ([]byte, error) {
	conn := self.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("HGET", self.buildKey(key), field))
}

func(self *RedisWrapper) HMGet(key string, fields []string) (map[string]string, error) {
	conn := self.Get()
	defer conn.Close()

	params := make([]interface{}, 0, len(fields)+1)
	params = append(params, self.buildKey(key))
	for _, v := range fields {
		params = append(params, v)
	}

	values, err := redis.Strings(conn.Do("HMGET",  params...))
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(fields))
	for i := 0; i < len(fields);i++ {
		result[fields[i]] = values[i]
	}

	return result, nil
}


//data = 1存在，data = 0不存在
func(self *RedisWrapper) HExist(key, field string)(int64, error) {
	conn := self.Get()
	defer conn.Close()

	data, err := redis.Int64(conn.Do("HEXISTS", self.buildKey(key), field))

	return data, err
}

func(self *RedisWrapper) HDel(key, field string) error {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", self.buildKey(key), field)

	return err
}

func(self *RedisWrapper) HGetAll(key string)(values []interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("HGETALL", self.buildKey(key)))
	return
}


func(self *RedisWrapper) HLen(key string)(size int, err error){
	conn := self.Get()
	defer conn.Close()

	size, err = redis.Int(conn.Do("HLEN", self.buildKey(key)))
	return
}

func(self *RedisWrapper) Del(key string) error {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", self.buildKey(key))

	return err
}

func(self *RedisWrapper) Expire(key string, seconds int64) error {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", self.buildKey(key), seconds)
	return err
}

/**
 * @param key string
 * @param seconds int64 unix时间戳，单位秒
 */
func(self *RedisWrapper) ExpireAt(key string, seconds int64) error {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIREAT", self.buildKey(key), seconds)
	return err
}


func(self *RedisWrapper) LPush(key string, value []byte) error {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", self.buildKey(key), value)

	return err
}

func(self *RedisWrapper)  RPop(key string) ([]byte, error) {

	conn := self.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("RPOP", self.buildKey(key)))
}

func(self *RedisWrapper)  FlushAll()(err error){
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHALL")
	return
}

func(self *RedisWrapper) FlushDB()(err error){
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHDB")
	return
}


func(self *RedisWrapper) SAdd(key string, member interface{}) (err error) {
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("SADD", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) SRem(key string, member interface{})(err error) {
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("SREM", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) SCard(key string)(size int, err error) {
	conn := self.Get()
	defer conn.Close()

	size, err = redis.Int(conn.Do("SREM", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SPop(key string)(value interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	value, err = redis.Int(conn.Do("SPOP", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SMembers(key string)(values []interface{}, err error) {
	conn := self.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("SMEMBERS", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SRandMember(key string, count int)(values []interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	values, err = redis.Values(conn.Do("SRANDMEMBER", self.buildKey(key), count))
	return
}

func(self *RedisWrapper) SIsMember(key string, member interface{})(value bool, err error){
	conn := self.Get()
	defer conn.Close()

	value, err = redis.Bool(conn.Do("SISMEMBER", self.buildKey(key), member))
	return
}


func(self *RedisWrapper) ZAdd(key string, score float64, value interface{})  error {
	conn := self.Get()
	defer conn.Close()

	var err error
	_, err = conn.Do("ZADD", self.buildKey(key), score, value)
	return err
}

func(self *RedisWrapper) ZCard(key string) (size int64, err error) {
	conn := self.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("ZCard", self.buildKey(key)))
}


//根据score获取数据
func(self *RedisWrapper) ZRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	conn := self.Get()
	defer conn.Close()

	if withScores {
		if count > 0 {
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", self.buildKey(key), min, max, "WITHSCORES", "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", self.buildKey(key), min, max, "WITHSCORES"))
		}
	}else{
		if count > 0 {
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", self.buildKey(key), min, max, "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZRANGEBYSCORE", self.buildKey(key), min, max))
		}
	}


	return
}

func(self *RedisWrapper) ZRevRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	conn := self.Get()
	defer conn.Close()

	if withScores {
		if count > 0 {
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", self.buildKey(key), min, max, "WITHSCORES", "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", self.buildKey(key), min, max, "WITHSCORES"))
		}
	}else{
		if count > 0 {
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", self.buildKey(key), min, max, "LIMIT", offset, count))
		}else{
			values, err = redis.Values(conn.Do("ZREVRANGEBYSCORE", self.buildKey(key), min, max))
		}
	}


	return
}

func(self *RedisWrapper) ZRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	conn := self.Get()
	defer conn.Close()

	if withScores {
		values, err = redis.Values(conn.Do("ZRANGE", self.buildKey(key), start, stop, "WITHSCORES"))
	}else{
		values, err = redis.Values(conn.Do("ZRANGE", self.buildKey(key), start, stop))
	}

	return
}

func(self *RedisWrapper) ZRevRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	conn := self.Get()
	defer conn.Close()

	if withScores {
		values, err = redis.Values(conn.Do("ZREVRANGE", self.buildKey(key), start, stop, "WITHSCORES"))
	}else{
		values, err = redis.Values(conn.Do("ZREVRANGE", self.buildKey(key), start, stop))
	}

	return
}

func(self *RedisWrapper) ZIncreBy(key string, increment float64, member interface{})(err error) {
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("ZINCRBY", self.buildKey(key), increment, member)
	return
}

//移除一个元素
func(self *RedisWrapper) ZRem(key string, member interface{})(err error) {
	conn := self.Get()
	defer conn.Close()

	_, err = conn.Do("ZREM", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) ZRank(key string, member interface{})(index int64, err error) {
	conn := self.Get()
	defer conn.Close()

	index, err = redis.Int64(conn.Do("ZRANK", self.buildKey(key), member))
	return
}

func(self *RedisWrapper) ZRevRank(key string, member interface{})(index int64, err error) {
	conn := self.Get()
	defer conn.Close()

	index, err = redis.Int64(conn.Do("ZREVRANK", self.buildKey(key), member))
	return
}

func(self *RedisWrapper) ZScore(key string, member interface{})(score float64, err error) {
	conn := self.Get()
	defer conn.Close()

	score, err = redis.Float64(conn.Do("ZSCORE", self.buildKey(key), member))
	return
}


func(self *RedisWrapper) SGet(key string) ([]byte, error) {

	conn := self.Get()
	defer conn.Close()

	return  redis.Bytes(conn.Do("GET", self.buildKey(key)))
}

func(self *RedisWrapper) SSet(key string, value []byte, ex int, px int, nx bool, xx bool) error {
	conn := self.Get()
	defer conn.Close()

	var err error

	if ex > 0 {
		if nx {
			_, err = conn.Do("SET", self.buildKey(key), value, "EX", ex, "NX")
		}else if xx {
			_, err = conn.Do("SET", self.buildKey(key), value, "EX", ex, "XX")
		}else{
			_, err = conn.Do("SET", self.buildKey(key), value, "EX", ex)
		}
	}else if px > 0 {
		if nx {
			_, err = conn.Do("SET", self.buildKey(key), value, "PX", px, "NX")
		}else if xx {
			_, err = conn.Do("SET", self.buildKey(key), value, "PX", px, "XX")
		}else{
			_, err = conn.Do("SET", self.buildKey(key), value, "PX", px)
		}
	}else{
		if nx {
			_, err = conn.Do("SET", self.buildKey(key), value, "NX")
		}else if xx {
			_, err = conn.Do("SET", self.buildKey(key), value, "XX")
		}else{
			_, err = conn.Do("SET", self.buildKey(key), value)
		}
	}

	return err
}

func(self *RedisWrapper) Incr(key string) (int64, error) {
	conn := self.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("INCR", self.buildKey(key)))
}

func(self *RedisWrapper) Exist(key string) (bool, error) {
	conn := self.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", self.buildKey(key)))
}