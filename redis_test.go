package redis_wrapper

import (
	"testing"
	"strconv"
	"fmt"
	"encoding/json"
	"time"
)


//goos: darwin
//goarch: amd64
//pkg: Languege/redis_wrapper
//BenchmarkRedis_Del-8     	   20000	     74581 ns/op
//BenchmarkRedis_Set-8     	   20000	     71803 ns/op
//BenchmarkRedis_HSet-8    	   20000	     77275 ns/op
//BenchmarkSRandMember-8   	   20000	     90124 ns/op
//PASS

func BenchmarkRedis_Del(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)
	for i := 0; i < b.N ; i++  {
		Del("test")
	}
}

func BenchmarkRedis_Set(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	var err error
	for i := 0; i < b.N ; i++  {
		err = Set("test", []byte("value"), 0, 0, false, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedis_HSet(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	var err error
	for i := 0; i < b.N ; i++  {
		_, err = HSet("hashkey", "field" + strconv.Itoa(i),[]byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedisWrapper_HGetAll(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)
	OpenTrace(10, 50)

	for i := 0; i< 500;i++ {
		HSet("hashkey", "field" + strconv.Itoa(i), []byte("2122"))
	}
	var(
		err error
	)
	for i := 0; i < b.N ; i++  {
		_, err = HGetAll("hashkey")
		if err != nil {
			b.Fatal(err)
		}
	}


	data, _ := json.Marshal(StatTraceInfo())

	fmt.Printf("%+v \n", string(data))
}

func BenchmarkSRandMember(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	for i := 0; i< 100;i++ {
		SAdd("skey", i)
	}

	var err error
	for i := 0; i < b.N ; i++  {
		_, err = SRandMember("skey",  10)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestSet(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	err := Set("test", []byte("value"), 60, 0, false, false)
	if err != nil {
		t.FailNow()
	}
}

func TestZAdd(t *testing.T) {

	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	err := ZAdd("zset_test", 1.00, []byte("value"))
	if err != nil {
		t.FailNow()
	}
}

func TestZRangeByScore(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	values, err := ZRangeByScore("zset_test", 0.0, 2.1, false, 0, 0)
	t.Log(values, err)
	if err != nil {
		t.FailNow()
	}
}


func TestExist(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)

	Del("test")
	ok, err := Exist("test")
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Fatalf("Should Be False")
	}
}


func TestRedisWrapper_HMSetValueWithExpire(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)


	ok, err := HMSetValue("hmsetlua", map[string]interface{}{
		"f1":"v1",
		"f2":"v2",
	}, 300)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ok)
}

func TestRedisWrapper_HSetValueWithExpire(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)


	ok, err := HSetValue("hmsetlua", "f1", "v1", 100)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ok)
}


func TestRedisWrapper_SafeTryLock(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC", 1, 1, time.Hour)


	releaseFunc, err := SafeTryLock("dlock", 10)
	if err != nil {
		t.Fatal(err)
	}

	defer releaseFunc()

	time.Sleep(time.Second)
}



