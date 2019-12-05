package redis_wrapper


/**
 *@author LanguageY++2013
 *2019/2/20 5:33 PM
 **/
func OpenTrace(tracePercentage int, options... interface{}) {
	wrapper.OpenTrace(tracePercentage, options...)
}


func StatTraceInfo() map[string]*CommandTrace {
	return wrapper.StatTraceInfo()
}