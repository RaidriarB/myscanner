package scan

import (
	"fmt"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
	"myscanner/settings"
	"strings"
	"time"
)

func ScanPorts(aliveHosts []string) types.TargetWithPorts {

	var NUM_OF_TASKS = settings.PORT_SCAN_THREADS
	var p = pool.NewPool(NUM_OF_TASKS)
	var upHostWithPorts = make(types.TargetWithPorts)

	//1. 定义端口存活性检测函数
	p.Function = func(i interface{}) interface{} {
		netloc := i.(string)
		if checkPortAlive(netloc) {
			slog.Debug(netloc, " is open")
			return netloc
		}
		return nil
	}

	//2. 把要检测的端口加入队列
	go func() {
		var SCANALL = settings.SCANALL

		var portlist []int

		switch settings.PORTLIST_MODE {
		case 1:
			portlist = settings.PORTLIST_SIMPLIFIED
		case 2:
			portlist = settings.PORTLIST_TOP1000
		case 3:
			portlist = settings.PORTLIST
		default:
			fmt.Println("PORTLIST_MODE的值不合法！ 默认使用top1000")
			portlist = settings.PORTLIST_TOP1000
		}

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
		}
	}()

	//4. 开始执行端口存活性探测任务
	p.Run()
	slog.Warning("端口存活性探测任务完成")
	return upHostWithPorts
}

func checkPortAlive(netloc string) bool {
	timeout := time.Duration(settings.PORT_SCAN_TIMEOUT)
	return gonmap.PortScan("tcp", netloc, timeout)
}
