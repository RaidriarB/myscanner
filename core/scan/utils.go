package scan

import (
	"errors"
	"math/rand"
	"myscanner/lib/gonmap"
	"myscanner/lib/httpfinger"
	"myscanner/lib/slog"
)

// 一列数，分成parts份，取第which份，求开始和结束的下标。（which从1开始计数）
func getBeginAndEnd(len int, parts int, which int) (int, int) {
	if len < 0 || parts <= 0 || which <= 0 || which > parts {
		panic(errors.New("getBeginAndEnd函数的参数不合法！"))
	}
	begin := (len * (which - 1)) / parts
	end := ((len * which) / parts) - 1
	//循环时只需要 begin <= i <= end
	return begin, end
}

//将一个字符串数组用特定的随机种子打乱
func shuffleStringArray(a []string, randID int64) {

	rand.Seed(randID)
	numOfHosts := len(a)
	rand.Shuffle(numOfHosts, func(i, j int) { a[i], a[j] = a[j], a[i] })
}

//初始化gonmap模块
func InitNmap() {
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条", r["FaviconHash"], r["KeywordFinger"])
	//gonmap探针/指纹库初始化
	r = gonmap.Init(9, 5*1000000000)
	slog.Infof("成功加载NMAP探针:[%d]个,指纹[%d]条", r["PROBE"], r["MATCH"])
}
