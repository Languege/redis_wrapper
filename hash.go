package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/
func HSet(key,field string, value []byte) (int64, error) {
	return wrapper.HSet(key, field, value)
}

func HGet(key,field string) ([]byte, error) {
	return wrapper.HGet(key, field)
}

//data = 1存在，data = 0不存在
func HExist(key, field string)(int64, error) {
	return wrapper.HExist(key, field)
}

func HDel(key, field string) error {
	return wrapper.HDel(key, field)
}

func HGetAll(key string)(values []interface{}, err error){
	return wrapper.HGetAll(key)
}

func HMSet(key string, kv map[string]string) (int64, error) {
	return wrapper.HMSet(key,  kv)
}

func HMGet(key string,fields []string) (map[string]string, error) {
	return wrapper.HMGet(key, fields)
}