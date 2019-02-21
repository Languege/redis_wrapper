package redis_wrapper

import (
	"testing"
	"strconv"
)


//goos: darwin
//goarch: amd64
//pkg: skddj/redis
//BenchmarkRedis_Del-8    	   20000	     79356 ns/op
//BenchmarkRedis_Set-8    	   20000	     77270 ns/op
//BenchmarkRedis_HSet-8   	   20000	     79877 ns/op
//PASS

func BenchmarkRedis_Del(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")
	for i := 0; i < b.N ; i++  {
		Del("test")
	}
}

func BenchmarkRedis_Set(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	var err error
	for i := 0; i < b.N ; i++  {
		err = Set("test", []byte("value"), 0, 0, false, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedis_HSet(b *testing.B) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	var err error
	for i := 0; i < b.N ; i++  {
		_, err = HSet("hashkey", "field" + strconv.Itoa(i),[]byte("value"))
		if err != nil {
			b.Fatal(err)
		}
	}
}


func TestSet(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	err := Set("test", []byte("value"), 60, 0, false, false)
	if err != nil {
		t.FailNow()
	}
}

func TestZAdd(t *testing.T) {

	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	err := ZAdd("zset_test", 1.00, []byte("value"))
	if err != nil {
		t.FailNow()
	}
}

func TestZRangeByScore(t *testing.T) {
	InitConnect("127.0.0.1", "6379", "SjhkHD3J5k6H8SjSbK3SC")

	values, err := ZRangeByScore("zset_test", 0.0, 2.1, false, 0, 0)
	t.Log(values, err)
	if err != nil {
		t.FailNow()
	}
}




