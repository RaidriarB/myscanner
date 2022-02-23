package main

import (
	"fmt"
	"myscanner/lib/gonmap"
	"myscanner/lib/httpfinger"
	"myscanner/lib/slog"
)

func add(a, b int) int {
	return a + b
}

func main() {
	//HTTP指纹库初始化
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条", r["FaviconHash"], r["KeywordFinger"])
	//gonmap探针/指纹库初始化
	r = gonmap.Init(9, 10000000000)
	slog.Infof("成功加载NMAP探针:[%d]个,指纹[%d]条", r["PROBE"], r["MATCH"])
	n := gonmap.New()
	//nl1 := "193.170.192.215:8080"
	nl2 := "127.0.0.1:80"
	m := gonmap.GetTcpBanner(nl2, n, 30000000000)
	fmt.Println(m)
}
