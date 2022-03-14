package scan

import (
	"fmt"
	"myscanner/lib/gonmap"
	"myscanner/lib/slog"
	"myscanner/settings"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

//根据config中的portlist模式来扫描
func ScanPorts(aliveHosts []string, SCANALL bool) sync.Map {
	var portlist []int
	switch settings.PORTLIST_MODE {
	case 1:
		portlist = settings.PORTLIST_EMPTY
	case 2:
		portlist = settings.PORTLIST_FOR_DEBUG
	case 3:
		portlist = settings.PORTLIST_SIMPLIFIED
	case 4:
		portlist = settings.PORTLIST_TOP1000
	case 5:
		portlist = settings.PORTLIST
	default:
		fmt.Println("PORTLIST_MODE的值不合法！ 默认使用top1000")
		portlist = settings.PORTLIST_TOP1000
	}
	return ScanPortsWithShuffle(aliveHosts, portlist, SCANALL, 1, 1, time.Now().UnixNano())
}

func ScanPortsWithShuffle(aliveHosts []string, portlist []int, SCANALL bool, parts int, which int, randID int64) sync.Map {

	var threads = settings.PORT_SCAN_THREADS

	var upHostWithPorts sync.Map
	var netlocsToScan []string

	var lock sync.Mutex
	var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(threads, func(i interface{}) {
		netloc := i.(string)
		if checkPortAlive(netloc) {
			slog.Debug(netloc, " is open")
			host := strings.Split(netloc, ":")[0]
			port := strings.Split(netloc, ":")[1]

			//千万别忘了加锁
			lock.Lock()
			lst := []string{}
			if value, ok := upHostWithPorts.Load(host); ok {
				lst = value.([]string)
			}
			lst = append(lst, port)
			upHostWithPorts.Store(host, lst)
			lock.Unlock()
		}
		wg.Done()
	})
	defer p.Release()

	//2. 把要检测的端口加入队列

	for _, host := range aliveHosts {
		for _, port := range portlist {
			netloc := host + ":" + fmt.Sprintf("%d", port)
			netlocsToScan = append(netlocsToScan, netloc)
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
				}
			}
		}
	}

	//3) 打乱并加入队列
	numOfNetlocs := len(netlocsToScan)

	fmt.Println("要扫描 ", numOfNetlocs, "个 net location")

	shuffleStringArray(netlocsToScan, randID)
	//fmt.Println("打乱后:", hostsToScan)

	begin, end := getBeginAndEnd(numOfNetlocs, parts, which)

	fmt.Printf("下标:%d->%d\n", begin, end)

	for i := begin; i <= end; i++ {
		wg.Add(1)
		fmt.Println("准备放入队列:", netlocsToScan[i])
		_ = p.Invoke(string(netlocsToScan[i]))
	}
	wg.Wait()
	fmt.Println("端口存活性探测任务完成")
	return upHostWithPorts
}

func checkPortAlive(netloc string) bool {

	fmt.Printf("Checking %s\n", netloc)
	if settings.DEV_MODE {
		time.Sleep(100 * time.Millisecond)
		return true
	} else {
		timeout := time.Duration(settings.PORT_SCAN_TIMEOUT)
		return gonmap.PortScan("tcp", netloc, timeout)
	}
}
