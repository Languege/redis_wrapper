package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/

func Get(key string) ([]byte, error) {

	return wrapper.SGet(key)
}

func Set(key string, value []byte, ex int, px int, nx bool, xx bool) error {
	return wrapper.SSet(key, value, ex, px, nx, xx)
}

func Incr(key string) (int64, error) {
	return wrapper.Incr(key)
}