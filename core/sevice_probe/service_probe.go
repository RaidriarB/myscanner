package service_probe

import (
	"fmt"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/httpfinger"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
)

func ServiceProbe(twp types.TargetWithPorts) types.TargetPortBanners {

	// TODO: pool的进程数也可以写进配置
	var p = pool.NewPool(10)
	var result = types.TargetPortBanners{}

	//TODO: 初始化代码需不需要调整一下
	r := httpfinger.Init()
	slog.Infof("成功加载favicon指纹:[%d]条，keyword指纹:[%d]条", r["FaviconHash"], r["KeywordFinger"])
	//gonmap探针/指纹库初始化
	r = gonmap.Init(9, 5*1000000000)
	slog.Infof("成功加载NMAP探针:[%d]个,指纹[%d]条", r["PROBE"], r["MATCH"])

	p.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		// 这个gonmap.New()不能用一个变量表示，否则会出现并发问题
		//var n = gonmap.New() 不可以！
		r := gonmap.GetTcpBanner(netloc, gonmap.New(), 20*1000000000)
		return r
	}

	//启用TCP层面协议识别任务下发器
	go func() {
		for host, ports := range twp {
			for _, port := range ports {
				netloc := host + ":" + port
				p.In <- netloc
			}
		}
		slog.Info("TCP层协议识别任务下发完毕")
		p.InDone()
	}()

	//启用TCP层指纹探测结果接受器
	go func() {
		for out := range p.Out {
			if out == nil {
				continue
			}
			tcpBanner := out.(*gonmap.TcpBanner)
			if tcpBanner == nil {
				continue
			}

			uri := tcpBanner.Target.URI()
			status := tcpBanner.Status
			service := tcpBanner.TcpFinger.Service
			slog.Debugf("%s %s %s", uri, status, service)
			result = append(result, tcpBanner)
		}
	}()

	//开始执行TCP层面协议识别任务
	p.Run()
	fmt.Println("TCP层协议识别任务完成")
	return result
}
