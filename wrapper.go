package redis_wrapper

import (
	"github.com/gomodule/redigo/redis"
	"time"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
	"sync"
)

const(
	Default_MaxStatNum = 1000
)
/**
 *@author LanguageY++2013
 *2019/8/9 9:12 PM
 **/
type RedisWrapper struct {
	redis.Pool
	Prefix 		string	//key前缀
	TracePercentage	int	//采集概率百分比
	TraceList		map[string]*CommandTrace //命令统计
	mutex 			sync.RWMutex	//统计锁
	MaxStatNum		int64	//最大统计次数，避免数据不能反映最近状况
}

type CommandTrace struct {
	TotalNum		int64
	TotalMs			int64
	AvgMsPerOp		int64
	ResetNum		int64	//因超出统计数次而重置次数
}

func NewRedisWrapper(ip string, port string, password string, maxIdle int, idleTimeout time.Duration, maxActive int) *RedisWrapper{
	addr := ip + ":" + port

	options := []redis.DialOption{
		redis.DialConnectTimeout(time.Second*20),
		redis.DialReadTimeout(time.Second*20),
		redis.DialWriteTimeout(time.Second*20),
	}

	if password != "" {
		options = append(options, redis.DialPassword(password))
	}
	w  :=  &RedisWrapper{
		Pool:redis.Pool{
			MaxIdle:     maxIdle,
			IdleTimeout: idleTimeout,

			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", addr, options...)
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			MaxActive: maxActive,
			Wait:      true,
		},
	}

	go w.closeConnection()

	return w
}

//开启追踪
func(self *RedisWrapper) OpenTrace(tracePercentage int, options... interface{}) {
	self.TracePercentage = tracePercentage
	if len(options) > 0 {
		if v, ok := options[0].(int);ok {
			self.MaxStatNum = int64(v)
		}
	}

	if self.MaxStatNum <= 0 {
		self.MaxStatNum = Default_MaxStatNum
	}


	self.mutex.Lock()
	self.TraceList = map[string]*CommandTrace{}
	self.mutex.Unlock()
}

func(self *RedisWrapper) Stat(command string, startTime time.Time) {
	if self.TracePercentage > 0 || self.TracePercentage < 100 {
		randIndex := rand.Intn(100)
		if randIndex < self.TracePercentage {
			//0~self.TracePercentage
			self.mutex.Lock()
			if v, ok := self.TraceList[command];ok {
				if v.TotalNum >= self.MaxStatNum {
					v.TotalNum = 1
					v.TotalMs = time.Now().Sub(startTime).Microseconds()
					v.ResetNum++
				}else{
					v.TotalNum++
					v.TotalMs += time.Now().Sub(startTime).Microseconds()
				}

			}else{
				self.TraceList[command] = &CommandTrace{TotalNum:1,TotalMs:time.Now().Sub(startTime).Microseconds()}
			}
			self.mutex.Unlock()
		}
	}
}

func(self *RedisWrapper) StatTraceInfo() map[string]*CommandTrace {
	result := map[string]*CommandTrace{}
	self.mutex.Lock()
	for k, v := range self.TraceList {
		if v.TotalNum > 0{
			self.TraceList[k].AvgMsPerOp = v.TotalMs / v.TotalNum
		}
	}
	self.mutex.Unlock()

	self.mutex.RLock()
	result = self.TraceList
	self.mutex.RUnlock()

	return result
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

func(self *RedisWrapper) SafeTryLock(key string, seconds int) (releaseCallback func() error, err error) {
	conn := self.Get()

	defer self.Stat("TryLock", time.Now())

	uniqueID := time.Now().UnixNano()
	_, err = redis.String(conn.Do("SET", self.buildKey(key), uniqueID, "EX", seconds, "NX"))
	if err != nil {
		return
	}

	releaseCallback = func() (err error){
				script := redis.NewScript(1, `
	 					if redis.call("get",KEYS[1]) == ARGV[1] then
          						return redis.call("del",KEYS[1])
      					else
          					return 0
      					end
					`)

				_, err = script.Do(conn, self.buildKey(key), uniqueID)

				conn.Close()
				return
	}

	return
}



func(self *RedisWrapper) TryLock(key string, seconds int)(uniqueID int64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("TryLock", time.Now())

	uniqueID = time.Now().UnixNano()
	_, err = redis.String(conn.Do("SET", self.buildKey(key), uniqueID, "EX", seconds, "NX"))
	return
}

func(self *RedisWrapper) Release(key string, uniqueID int64)(err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("Release", time.Now())

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

	defer self.Stat("HSET", time.Now())

	return redis.Int64(conn.Do("HSET", self.buildKey(key), field, value))
}

func(self *RedisWrapper) HSetValue(key,field string, value interface{}, extra... int64) (int64, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HSET", time.Now())

	if len(extra) > 0 && extra[0] > 0 {
		params := []interface{}{self.buildKey(key),field,value,extra[0]}
		script := redis.NewScript(1, `
				local ret = redis.call('HSET', KEYS[1], ARGV[1], ARGV[2])
				redis.call('expire', KEYS[1], ARGV[3])
	 			return ret
		`)

		return redis.Int64(script.Do(conn, params...))
	}

	return redis.Int64(conn.Do("HSET", self.buildKey(key), field, value))
}


//执行成功返回OK
func(self *RedisWrapper) HMSet(key string, kv map[string]string) (string, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HMSET", time.Now())

	params := make([]interface{}, 0, 2 * len(kv)+1)
	params = append(params, self.buildKey(key))
	for k, v := range kv  {
		params = append(params, k, v)
	}

	return redis.String(conn.Do("HMSET",  params...))
}

//执行成功返回OK
func(self *RedisWrapper) HMSetValue(key string, kv map[string]interface{}, extra... int64) (string, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HMSET", time.Now())

	params := make([]interface{}, 0, 2 * len(kv)+1)
	params = append(params, self.buildKey(key))
	for k, v := range kv  {
		params = append(params, k, v)
	}

	if len(extra) > 0 && extra[0] > 0 {
		params = append(params, extra[0])
		script := redis.NewScript(1, `
	 local mi = table.maxn(ARGV)
	 local expireSeconds = ARGV[mi]
	 table.remove(ARGV, mi)
	 local ret = redis.call('HMSET', KEYS[1], unpack(ARGV))
	 redis.call('expire', KEYS[1], expireSeconds)
	 return ret
	`)
		return redis.String(script.Do(conn, params...))
	}

	return redis.String(conn.Do("HMSET",  params...))
}

func(self *RedisWrapper) HGet(key,field string) ([]byte, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HGET", time.Now())

	return  redis.Bytes(conn.Do("HGET", self.buildKey(key), field))
}

func(self *RedisWrapper) HGetInt64(key,field string) (int64, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HGET", time.Now())

	return  redis.Int64(conn.Do("HGET", self.buildKey(key), field))
}

func(self *RedisWrapper) HGetString(key,field string) (string, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HGET", time.Now())

	return  redis.String(conn.Do("HGET", self.buildKey(key), field))
}

func(self *RedisWrapper) HMGet(key string, fields []string) (map[string]string, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HMGET", time.Now())

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

func(self *RedisWrapper) HMGetInt(key string, fields []string) (map[string]int, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HMGET", time.Now())

	params := make([]interface{}, 0, len(fields)+1)
	params = append(params, self.buildKey(key))
	for _, v := range fields {
		params = append(params, v)
	}

	values, err := redis.Ints(conn.Do("HMGET",  params...))
	if err != nil {
		return nil, err
	}

	result := make(map[string]int, len(fields))
	for i := 0; i < len(fields);i++ {
		result[fields[i]] = values[i]
	}

	return result, nil
}


//data = 1存在，data = 0不存在
func(self *RedisWrapper) HExist(key, field string)(int64, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HEXISTS", time.Now())

	data, err := redis.Int64(conn.Do("HEXISTS", self.buildKey(key), field))

	return data, err
}

func(self *RedisWrapper) HDel(key, field string) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HDEL", time.Now())

	_, err := conn.Do("HDEL", self.buildKey(key), field)

	return err
}

func(self *RedisWrapper) HGetAll(key string)(values []interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HGETALL", time.Now())

	values, err = redis.Values(conn.Do("HGETALL", self.buildKey(key)))
	return
}

func(self *RedisWrapper) HGetAllInt(key string)(values map[string]int, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HGETALL", time.Now())

	values, err = redis.IntMap(conn.Do("HGETALL", self.buildKey(key)))
	return
}

func(self *RedisWrapper) HGetAll2Map(key string)(ret map[string][]byte, err error) {
	ret = map[string][]byte{}
	var values []interface{}
	values, err =  self.HGetAll(key)
	if err != nil {
		return
	}

	if len(values) == 0 {
		return
	}

	for i := 0; i < len(values); i += 2 {
		key := string(values[i].([]byte))
		data := values[i+1].([]byte)

		ret[key] = data
	}
	return
}


func(self *RedisWrapper) HLen(key string)(size int, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HLEN", time.Now())

	size, err = redis.Int(conn.Do("HLEN", self.buildKey(key)))
	return
}


func(self *RedisWrapper) HIncrBy(key, field string, increment int64)(ret int64, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("HINCRBY", time.Now())

	ret, err = redis.Int64(conn.Do("HINCRBY", self.buildKey(key), field, increment))
	return
}



func(self *RedisWrapper) Del(key string) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("DEL", time.Now())

	_, err := conn.Do("DEL", self.buildKey(key))

	return err
}

func(self *RedisWrapper) Expire(key string, seconds int) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("Expire", time.Now())

	_, err := conn.Do("EXPIRE", self.buildKey(key), seconds)
	return err
}

/**
 * @param key string
 * @param seconds int64 unix时间戳，单位秒
 */
func(self *RedisWrapper) ExpireAt(key string, seconds int) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("EXPIREAT", time.Now())

	_, err := conn.Do("EXPIREAT", self.buildKey(key), seconds)
	return err
}


func(self *RedisWrapper) LPush(key string, value []byte) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("LPUSH", time.Now())

	_, err := conn.Do("LPUSH", self.buildKey(key), value)

	return err
}

func(self *RedisWrapper) RPush(key string, value []byte) error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("RPUSH", time.Now())

	_, err := conn.Do("RPUSH", self.buildKey(key), value)

	return err
}

func(self *RedisWrapper)  RPop(key string) ([]byte, error) {

	conn := self.Get()
	defer conn.Close()

	defer self.Stat("RPOP", time.Now())

	return  redis.Bytes(conn.Do("RPOP", self.buildKey(key)))
}

func(self *RedisWrapper)  LPop(key string) ([]byte, error) {

	conn := self.Get()
	defer conn.Close()

	defer self.Stat("LPOP", time.Now())

	return  redis.Bytes(conn.Do("LPOP", self.buildKey(key)))
}

func(self *RedisWrapper) LRange(key string, start, stop int)(ret []string, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("LRANGE", time.Now())

	return redis.Strings(conn.Do("LRANGE", self.buildKey(key), start, stop))
}

func(self *RedisWrapper) LLen(key string)(ret int64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("LLEN", time.Now())

	return redis.Int64(conn.Do("LLEN", self.buildKey(key)))
}

func(self *RedisWrapper)  FlushAll()(err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("FLUSHALL", time.Now())

	_, err = conn.Do("FLUSHALL")
	return
}

func(self *RedisWrapper) FlushDB()(err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("FLUSHDB", time.Now())

	_, err = conn.Do("FLUSHDB")
	return
}


func(self *RedisWrapper) SAdd(key string, member interface{}) (err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SADD", time.Now())

	_, err = conn.Do("SADD", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) SRem(key string, member interface{})(err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SREM", time.Now())

	_, err = conn.Do("SREM", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) SCard(key string)(size int, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SCARD", time.Now())

	size, err = redis.Int(conn.Do("SCARD", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SPop(key string)(value interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SPOP", time.Now())

	value, err = redis.Int(conn.Do("SPOP", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SMembers(key string)(values []interface{}, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SMEMBERS", time.Now())

	values, err = redis.Values(conn.Do("SMEMBERS", self.buildKey(key)))
	return
}

func(self *RedisWrapper) SRandMember(key string, count int)(values []interface{}, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SRANDMEMBER", time.Now())

	values, err = redis.Values(conn.Do("SRANDMEMBER", self.buildKey(key), count))
	return
}

func(self *RedisWrapper) SIsMember(key string, member interface{})(value bool, err error){
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("SISMEMBER", time.Now())

	value, err = redis.Bool(conn.Do("SISMEMBER", self.buildKey(key), member))
	return
}


func(self *RedisWrapper) ZAdd(key string, score float64, value interface{})  error {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZADD", time.Now())

	var err error
	_, err = conn.Do("ZADD", self.buildKey(key), score, value)
	return err
}

func(self *RedisWrapper) ZCard(key string) (size int64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZCard", time.Now())

	return redis.Int64(conn.Do("ZCard", self.buildKey(key)))
}


//根据score获取数据
func(self *RedisWrapper) ZRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZRANGEBYSCORE", time.Now())

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

	defer self.Stat("ZREVRANGEBYSCORE", time.Now())

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

	defer self.Stat("ZRANGE", time.Now())

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

	defer self.Stat("ZREVRANGE", time.Now())

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

	defer self.Stat("ZINCRBY", time.Now())

	_, err = conn.Do("ZINCRBY", self.buildKey(key), increment, member)
	return
}

//移除一个元素
func(self *RedisWrapper) ZRem(key string, member interface{})(err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZREM", time.Now())

	_, err = conn.Do("ZREM", self.buildKey(key), member)
	return
}

func(self *RedisWrapper) ZRank(key string, member interface{})(index int64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZRANK", time.Now())

	index, err = redis.Int64(conn.Do("ZRANK", self.buildKey(key), member))
	return
}

func(self *RedisWrapper) ZRevRank(key string, member interface{})(index int64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZREVRANK", time.Now())

	index, err = redis.Int64(conn.Do("ZREVRANK", self.buildKey(key), member))
	return
}

func(self *RedisWrapper) ZScore(key string, member interface{})(score float64, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZSCORE", time.Now())

	score, err = redis.Float64(conn.Do("ZSCORE", self.buildKey(key), member))
	return
}

func(self *RedisWrapper) ZRemRangeByScore(key string, min, max float64)(num int, err error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("ZREMRANGEBYSCORE", time.Now())

	num, err = redis.Int(conn.Do("ZREMRANGEBYSCORE", self.buildKey(key), min, max))
	return
}


func(self *RedisWrapper) SGet(key string) ([]byte, error) {

	conn := self.Get()
	defer conn.Close()

	defer self.Stat("GET", time.Now())

	return  redis.Bytes(conn.Do("GET", self.buildKey(key)))
}

func(self *RedisWrapper) SGetInt64(key string) (int64, error) {

	conn := self.Get()
	defer conn.Close()

	defer self.Stat("GET", time.Now())

	return  redis.Int64(conn.Do("GET", self.buildKey(key)))
}

func(self *RedisWrapper) SGetString(key string) (string, error) {

	conn := self.Get()
	defer conn.Close()

	defer self.Stat("GET", time.Now())

	return  redis.String(conn.Do("GET", self.buildKey(key)))
}

func(self *RedisWrapper) SSet(key string, value []byte, ex int, px int, nx bool, xx bool) error {
	return self.SSetValue(key, value, ex, px, nx, xx)
}

func(self *RedisWrapper) SSetValue(key string, value interface{}, options... interface{}) error {
	var(
		ex, px int
		nx, xx bool
	)
	if len(options) > 0 {
		if v, ok := options[0].(int);ok {
			ex = v
		}
	}
	if len(options) > 1 {
		if v, ok := options[1].(int);ok {
			px = v
		}
	}
	if len(options) > 2 {
		if v, ok := options[2].(bool);ok {
			nx = v
		}
	}
	if len(options) > 3 {
		if v, ok := options[3].(bool);ok {
			xx = v
		}
	}
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("GET", time.Now())

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

	defer self.Stat("INCR", time.Now())

	return redis.Int64(conn.Do("INCR", self.buildKey(key)))
}

func(self *RedisWrapper) Exist(key string) (bool, error) {
	conn := self.Get()
	defer conn.Close()

	defer self.Stat("EXISTS", time.Now())

	return redis.Bool(conn.Do("EXISTS", self.buildKey(key)))
}

func(self *RedisWrapper) TTL(key string) (int64, error){
	conn := self.Get()
	defer  conn.Close()

	defer self.Stat("TTL", time.Now())

	return redis.Int64(conn.Do("TTL", self.buildKey(key)))
}

func(self *RedisWrapper) IncrBy(key string, increment int) (int, error) {
	conn := self.Get()
	defer conn.Close()

	return redis.Int(conn.Do("INCRBY", self.buildKey(key), increment))
}