package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/22 6:32 PM
 **/
func FlushAll()(err error){
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHALL")
	return
}

func FlushDB()(err error){
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("FLUSHDB")
	return
}

