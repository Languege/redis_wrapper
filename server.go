package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/22 6:32 PM
 **/
func FlushAll()(err error){
	return wrapper.FlushAll()
}

func FlushDB()(err error){
	return wrapper.FlushDB()
}

