package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 6:00 PM
 **/
func Del(key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)

	return err
}

func Expire(key string, seconds int64) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", key, seconds)
	return err
}

/**
 * @param key string
 * @param seconds int64 unix时间戳，单位秒
 */
func ExpireAt(key string, seconds int64) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIREAT", key, seconds)
	return err
}