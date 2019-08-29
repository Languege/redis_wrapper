package main

import (
	"github.com/Languege/redis_wrapper"
	"fmt"
)

/**
 *@author LanguageY++2013
 *2019/2/20 7:51 PM
 **/
func main(){

	redis_wrapper.InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

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

	//for i := 0; i< 10;i++ {
	//	redis_wrapper.HSet("hashtable", "field" + strconv.Itoa(i), []byte(encoding.JSON("field" + strconv.Itoa(i))))
	//}
	//
	//values, err = redis_wrapper.HGetAll("hashtable")
	//if err == nil {
	//	for i := 0; i < len(values); i+=2  {
	//		fmt.Println(string(values[i].([]byte)), string(values[i+1].([]byte)))
	//	}
	//}

	//HMSet, HMGet
	n, err := redis_wrapper.HMSet("hmset_test", map[string]string{"f1":"v1","f2":"v2"})
	fmt.Println(n)
	fmt.Println(err)
	result, err := redis_wrapper.HMGet("hmset_test", []string{"f1","f2"})
	fmt.Println(result)
	fmt.Println(err)


	redis_wrapper.ZAdd("zadd_test", 1.1, 1000)

	score, err := redis_wrapper.ZScore("zadd_test", 1000)

	fmt.Println(score)
	fmt.Println(err)

}