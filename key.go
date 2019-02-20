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