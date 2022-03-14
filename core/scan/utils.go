package scan

import "errors"

// 一列数，分成parts份，取第which份，求开始和结束的下标
// which从0开始计数
func getBeginAndEnd(len int, parts int, which int) (int, int) {
	if len < 0 || parts <= 0 || which <= 0 || which > parts {
		panic(errors.New("getBeginAndEnd函数的参数不合法！"))
	}
	begin := (len * (which - 1)) / parts
	end := ((len * which) / parts) - 1
	//循环时只需要 begin <= i <= end
	return begin, end
}
