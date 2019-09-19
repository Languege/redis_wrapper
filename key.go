package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 6:00 PM
 **/
func Del(key string) error {
	return wrapper.Del(key)
}

func Expire(key string, seconds int64) error {
	return wrapper.Expire(key, seconds)
}

/**
 * @param key string
 * @param seconds int64 unix时间戳，单位秒
 */
func ExpireAt(key string, seconds int64) error {
	return wrapper.ExpireAt(key, seconds)
}

func Exist(key string) (bool, error) {
	return wrapper.Exist(key)
}