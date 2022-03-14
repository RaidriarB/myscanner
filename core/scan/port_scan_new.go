package scan

import (
	"fmt"
	"math/rand"
	"myscanner/lib/slog"
	"myscanner/settings"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
)

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

	rand.Seed(randID)
	rand.Shuffle(numOfNetlocs, func(i, j int) { netlocsToScan[i], netlocsToScan[j] = netlocsToScan[j], netlocsToScan[i] })

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
