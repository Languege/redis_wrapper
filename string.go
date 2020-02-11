package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/

func Get(key string) ([]byte, error) {
	return wrapper.SGet(key)
}

func GetInt64(key string) (int64, error) {
	return wrapper.SGetInt64(key)
}


func GetString(key string) (string, error) {
	return wrapper.SGetString(key)
}

func Set(key string, value []byte, ex int, px int, nx bool, xx bool) error {
	return wrapper.SSet(key, value, ex, px, nx, xx)
}

func SetValue(key string, value interface{}, options... interface{}) error {
	return wrapper.SSetValue(key, value, options)
}

func Incr(key string) (int64, error) {
	return wrapper.Incr(key)
}

func IncrBy(key string, increment int)(int, error) {
	return wrapper.IncrBy(key, increment)
}