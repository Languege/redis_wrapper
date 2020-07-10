package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/22 6:36 PM
 **/
//分布式锁
func TryLock(key string, seconds int)(uniqueID int64, err error) {
	return wrapper.TryLock(key, seconds)
}

func Release(key string, uniqueID int64)(err error) {
	return wrapper.Release(key, uniqueID)
}


func SafeTryLock(key string, seconds int) (releaseFunc func(), err error) {
	return wrapper.SafeTryLock(key, seconds)
}

