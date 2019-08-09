package redis_wrapper

/**
 *@author LanguageY++2013
 *2019/2/20 5:32 PM
 **/

func ZAdd(key string, score float64, value interface{})  error {
	return wrapper.ZAdd(key, score, value)
}

func ZCard(key string) (size int64, err error) {
	return wrapper.ZCard(key)
}


//根据score获取数据
func ZRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	return wrapper.ZRangeByScore(key, min, max, withScores, offset, count)
}

func ZRevRangeByScore(key string, min float64, max float64, withScores bool, offset int, count int)(values []interface{}, err error)  {
	return wrapper.ZRevRangeByScore(key, min, max, withScores, offset, count)
}

func ZRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	return wrapper.ZRange(key, start, stop, withScores)
}

func ZRevRange(key string, start int, stop int, withScores bool) (values []interface{}, err error) {
	return wrapper.ZRevRange(key, start, stop, withScores)
}

func ZIncreBy(key string, increment float64, member interface{})(err error) {
	return wrapper.ZIncreBy(key, increment, member)
}

//移除一个元素
func ZRem(key string, member interface{})(err error) {
	return wrapper.ZRem(key, member)
}

func ZRank(key string, member interface{})(index int64, err error) {
	return wrapper.ZRank(key, member)
}

func ZRevRank(key string, member interface{})(index int64, err error) {
	return wrapper.ZRevRank(key, member)
}
