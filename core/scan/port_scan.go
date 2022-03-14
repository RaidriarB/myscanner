package scan

import (
	"fmt"
	"math/rand"
	"myscanner/lib/gonmap"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
	"myscanner/settings"
	"strings"
	"sync"
	"time"
)

func ScanPorts(aliveHosts []string, SCANALL bool) sync.Map {
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
	return ScanPortsWithShuffle(aliveHosts, portlist, SCANALL, 1, 1, time.Now().UnixNano())
}

func ScanPortsWithShuffle_old(aliveHosts []string, portlist []int, SCANALL bool, parts int, which int, randID int64) sync.Map {

	var NUM_OF_TASKS = settings.PORT_SCAN_THREADS
	var p = pool.NewPool(NUM_OF_TASKS)
	var upHostWithPorts sync.Map
	var netlocsToScan []string

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

		for _, host := range aliveHosts {
			for _, port := range portlist {
				netloc := host + ":" + fmt.Sprintf("%d", port)
				netlocsToScan = append(netlocsToScan, netloc)
				//p.In <- netloc
			}
			if SCANALL {

				isScanned := [settings.MAXPORT]bool{}
				for _, port := range portlist {
					isScanned[int(port)] = true
				}
				for port := 1; port < settings.MAXPORT; port++ {
					if !isScanned[port] {
						netloc := host + ":" + fmt.Sprintf("%d", port)
						netlocsToScan = append(netlocsToScan, netloc)
						//p.In <- netloc
					}
				}

			}
		}

		//3) 打乱并加入队列
		numOfNetlocs := len(netlocsToScan)

		fmt.Println("要扫描 ", numOfNetlocs, "个 net location")

		rand.Seed(randID)
		rand.Shuffle(numOfNetlocs, func(i, j int) { netlocsToScan[i], netlocsToScan[j] = netlocsToScan[j], netlocsToScan[i] })

		//fmt.Println("打乱后:", hostsToScan)

		begin, end := getBeginAndEnd(numOfNetlocs, parts, which)

		fmt.Printf("下标:%d->%d\n", begin, end)

		for i := begin; i <= end; i++ {
			fmt.Println("准备放入队列:", netlocsToScan[i])
			p.In <- netlocsToScan[i]
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

			lst := []string{}
			if value, ok := upHostWithPorts.Load(host); ok {
				lst = value.([]string)
			}
			lst = append(lst, port)
			upHostWithPorts.Store(host, lst)
		}
	}()

	//4. 开始执行端口存活性探测任务
	p.Run()
	slog.Warning("端口存活性探测任务完成")
	return upHostWithPorts
}

func checkPortAlive(netloc string) bool {

	fmt.Printf("Checking %s\n", netloc)
	if settings.DEV_MODE {
		time.Sleep(300 * time.Millisecond)
		return true
	} else {
		timeout := time.Duration(settings.PORT_SCAN_TIMEOUT)
		return gonmap.PortScan("tcp", netloc, timeout)
	}
}
