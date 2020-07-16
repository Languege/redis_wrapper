package main

import (
	"github.com/Languege/redis_wrapper"
	"fmt"
	"strconv"
	"time"
	"flag"
	"log"
)

var(
	host, port, password string
)

func init() {
	flag.StringVar(&host, "h", "", "-h")
	flag.StringVar(&port, "p", "", "-p")
	flag.StringVar(&password, "a", "", "-a")
	flag.Parse()
}
/**
 *@author LanguageY++2013
 *2019/2/20 7:51 PM
 **/
func main(){
	redis_wrapper.InitConnect(host, port, password, 0, 0, time.Hour)

	//var err error
	//values, err := redis_wrapper.ZRangeByScore("zset_test", 0.0, 2.1, true, 0,0)
	//
	//fmt.Println(values,err)
	//
	//for k, v := range values {
	//	fmt.Println(k, string(v.([]byte)))
	//}

	//uniqueID, err := redis_wrapper.TryLock("dlock", 1000)
	//time.Sleep(time.Duration(1) * time.Minute)
	//if err == nil {
	//	redis_wrapper.Release("dlock", uniqueID)
	//}

	//发布订阅测试
	//pubSubConn := redis.PubSubConn{Conn:redis_wrapper.GetConn()}
	//
	//err = pubSubConn.Subscribe("test_channel")
	//if err == nil {
	//	data := pubSubConn.ReceiveWithTimeout(time.Second)
	//	fmt.Println(data)
	//
	//	for {
	//		data = pubSubConn.Receive()
	//		msg, ok :=  data.(redis.Message)
	//		if ok {
	//			fmt.Println(msg, ok)
	//		}
	//	}
	//}

	//集合随机数测试
	//for i := 0; i< 100;i++ {
	//	redis_wrapper.SAdd("skey", i)
	//}
	//
	//ml, err := redis_wrapper.SRandMember("skey", 10)
	//if err == nil {
	//	for _, v := range ml {
	//		member, _ := strconv.ParseInt(string(v.([]byte)), 10, 64)
	//		fmt.Println(member)
	//	}
	//}

	for i := 0; i< 500;i++ {
		redis_wrapper.HSet("hashtable", "field" + strconv.Itoa(i), []byte("2122"))
	}
	//
	//values, err := redis_wrapper.HGetAll2Map("hashtable")
	//if err == nil {
	//	for k, v := range values {
	//		fmt.Printf("key:%s, value:%s \n", k, string(v))
	//	}
	//
	//}

	//HMSet, HMGet
	//n, err := redis_wrapper.HMSet("hmset_test", map[string]string{"f1":"v1","f2":"v2"})
	//fmt.Println(n)
	//fmt.Println(err)
	//result, err := redis_wrapper.HMGet("hmset_test", []string{"f1","f2"})
	//fmt.Println(result)
	//fmt.Println(err)
	//
	//size, err := redis_wrapper.HLen("hmset_tes")
	//fmt.Println(size)
	//fmt.Println(err)
	//
	//
	//redis_wrapper.ZAdd("zadd_test", 1.1, 1000)
	//
	//score, err := redis_wrapper.ZScore("zadd_test", 1000)
	//
	//fmt.Println(score)
	//fmt.Println(err)
	//
	//
	//ret, err := redis_wrapper.Exist("zadd_test1")
	//fmt.Println(ret)
	//fmt.Println(err)


	//生命周期检测
	//r1, err := redis_wrapper.TTL("no_exist_ttl_test")
	//fmt.Println(r1)
	//fmt.Println(err)
	//
	//redis_wrapper.Set("exist_no_ttl_test", []byte("exist_no_ttl_test"), 0, 0, false, false)
	//r2, err := redis_wrapper.TTL("exist_no_ttl_test")
	//fmt.Println(r2)
	//fmt.Println(err)
	//
	//redis_wrapper.Set("exist_ttl_test", []byte("exist_ttl_test"), 100, 0, false, false)
	//time.Sleep(time.Duration(10) * time.Second)
	//r3, err := redis_wrapper.TTL("exist_ttl_test")
	//fmt.Println(r3)
	//fmt.Println(err)

	//lkey := "lrangekey"
	//for i := 0; i< 10;i++ {
	//	redis_wrapper.RPush(lkey, []byte("field" + strconv.Itoa(i)))
	//}
	//
	//ret, err := redis_wrapper.LRange(lkey, 0, -1)
	//fmt.Println(ret)
	//fmt.Println(err)

	lkey := "zset_test"
	for i := 0; i< 10;i++ {
		redis_wrapper.ZAdd(lkey, float64(i), []byte("field" + strconv.Itoa(i)))
	}

	ret, err := redis_wrapper.ZRemRangeByScore(lkey, 0, 99999)
	fmt.Println(ret)
	fmt.Println(err)

	redis_wrapper.Set("string_test", []byte("1"), 0, 0, false, false)
	tmp, err := redis_wrapper.Get("string_test")
	fmt.Println(tmp)
	fmt.Println(err)


	ret, err = strconv.Atoi(string(tmp))
	fmt.Println(ret)
	fmt.Println(err)


	ret2, err2 := redis_wrapper.HIncrBy("hash_test", "1", 1)
	fmt.Println(ret2)
	fmt.Println(err2)


	ret3, err3 := redis_wrapper.HGetInt64("hash_test", "1")
	fmt.Println(ret3)
	fmt.Println(err3)


	redis_wrapper.SetValue("test_string", 100)
	ret4, err4 := redis_wrapper.GetInt64("test_string")
	fmt.Println(ret4)
	fmt.Println(err4)


	redis_wrapper.HSetValue("hluatest", "k1", "v1", 1000)

	data, err := redis_wrapper.HGet("hluatest", "k1")
	log.Println(err)
	log.Println(string(data))
}