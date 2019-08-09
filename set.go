package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:32 PM
 **/
func SAdd(key string, member interface{}) (err error) {
	return wrapper.SAdd(key, member)
}

func SRem(key string, member interface{})(err error) {
	return wrapper.SRem(key, member)
}

func SCard(key string)(size int, err error) {
	return wrapper.SCard(key)
}

func SPop(key string)(value interface{}, err error){
	return wrapper.SPop(key)
}

func SMembers(key string)(values []interface{}, err error) {
	return wrapper.SMembers(key)
}

func SRandMember(key string, count int)(values []interface{}, err error){
	return wrapper.SRandMember(key, count)
}

func SIsMember(key string, member interface{})(value bool, err error){
	return wrapper.SIsMember(key, member)
}
