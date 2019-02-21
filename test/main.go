package main

import (
	"Languege/redis_wrapper"
	"fmt"
)

/**
 *@author LanguageY++2013
 *2019/2/20 7:51 PM
 **/
func main(){

	redis_wrapper.InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	values, err := redis_wrapper.ZRangeByScore("zset_test", 0.0, 2.1, true, 0,0)

	fmt.Println(values,err)

	for k, v := range values {
		fmt.Println(k, string(v.([]byte)))
	}

}