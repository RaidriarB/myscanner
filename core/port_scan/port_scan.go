package port_scan

import (
	"fmt"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
	"strings"
)

func ScanPorts(aliveHosts []string) types.TargetWithPorts {

	var p = pool.NewPool(10)
	var upHostWithPorts = make(types.TargetWithPorts)

	//1. 定义端口存活性检测函数
	p.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		// TODO: timeout（3）变成配置
		if checkPortAlive(netloc) {
			slog.Debug(netloc, " is open")
			return netloc
		}
		return nil
	}

	//2. 把要检测的端口加入队列
	go func() {
		// TODO: 添加一个配置变量 SCANALL ，代表扫描所有端口
		var SCANALL = false

		// TODO: 把portlist移动到配置文件中
		portlist := []int{20, 21, 22, 52, 80, 443, 3306, 8080, 11111}

		for _, host := range aliveHosts {
			for _, port := range portlist {
				netloc := host + ":" + fmt.Sprintf("%d", port)
				p.In <- netloc
			}
			if SCANALL {
				isScanned := [65536]bool{}
				for _, port := range portlist {
					isScanned[int(port)] = true
				}
				for port := 1; port <= 65535; port++ {
					if !isScanned[port] {
						netloc := host + ":" + fmt.Sprintf("%d", port)
						p.In <- netloc
					}
				}

			}
		}
		slog.Info("端口存活性探测任务下发完毕")
		p.InDone()
	}()

	//启用端口存活性探测结果接受器
	go func() {
		for out := range p.Out {
			netloc := out.(string)
			host := strings.Split(netloc, ":")[0]
			port := strings.Split(netloc, ":")[1]

			upHostWithPorts[host] = append(upHostWithPorts[host], port)

			// FIXME: 下面这些是watchDog的代码
			// if host.IsOpenPort() == false && host.Length() == len(k.config.Port) && k.config.ClosePing == false {
			// 	url := fmt.Sprintf("icmp://%s", host.addr)
			// 	description := color.Red(color.Overturn("Not Open Any Port"))
			// 	output := fmt.Sprintf("%-30v %-26v %s", url, "Up", description)
			// 	k.watchDog.output <- output
			// }
			// 	upHosts.Set(host.addr, host)
			// }
		}
	}()

	//4. 开始执行端口存活性探测任务
	p.Run()
	slog.Warning("端口存活性探测任务完成")
	return upHostWithPorts
}

func checkPortAlive(netloc string) bool {

	// if strings.HasSuffix(netloc, "1") || strings.HasSuffix(netloc, "3") || strings.HasSuffix(netloc, "5") || strings.HasSuffix(netloc, "7") || strings.HasSuffix(netloc, "9") {
	// 	return true
	// }
	return gonmap.PortScan("tcp", netloc, 5*1000000000)
}
