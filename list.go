package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:31 PM
 **/
func LPush(key string, value []byte) error {
	return wrapper.LPush(key, value)
}

func RPush(key string, value []byte) error {
	return wrapper.RPush(key, value)
}

func RPop(key string) ([]byte, error) {
	return wrapper.RPop(key)
}

func LPop(key string) ([]byte, error) {
	return wrapper.LPop(key)
}

func LRange(key string, start, stop int)(ret []string, err error) {
	return wrapper.LRange(key, start, stop)
}

func LLen(key string)(ret int64, err error) {
	return wrapper.LLen(key)
}