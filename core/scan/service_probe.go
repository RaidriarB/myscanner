package scan

import (
	"fmt"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/slog"
	"myscanner/settings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

func ServiceProbe(twp sync.Map) types.TargetPortBanners {
	return ServiceProbeWithShuffle(twp, 1, 1, time.Now().UnixNano())
}

func ServiceProbeWithShuffle(twp sync.Map, parts int, which int, randID int64) types.TargetPortBanners {

	var wg sync.WaitGroup
	var result = types.TargetPortBanners{}
	var threads = settings.SERVICE_PROBE_THREADS

	init_nmap()

	var lock sync.Mutex
	p, _ := ants.NewPoolWithFunc(threads, func(i interface{}) {
		netloc := i.(string)
		// 这个gonmap.New()不能用一个变量表示，否则会出现并发问题
		//var n = gonmap.New() 不可以！
		tcpBanner := gonmap.GetTcpBanner(netloc, gonmap.New(), 20*1000000000)

		if tcpBanner != nil {
			uri := tcpBanner.Target.URI()
			status := tcpBanner.Status
			service := tcpBanner.TcpFinger.Service
			slog.Debugf("%s %s %s", uri, status, service)
			lock.Lock()
			result = append(result, tcpBanner)
			lock.Unlock()
		}
		wg.Done()
	})
	defer p.Release()

	var netlocsToProbe []string
	twp.Range(func(host, ports interface{}) bool {
		for _, port := range ports.([]string) {
			netloc := host.(string) + ":" + port
			netlocsToProbe = append(netlocsToProbe, netloc)
		}
		return true
	})
	numOfProbes := len(netlocsToProbe)
	fmt.Println("要服务识别 ", numOfProbes, "个 loc")

	shuffleStringArray(netlocsToProbe, randID)
	//fmt.Println("打乱后:", netlocsToProbe)

	begin, end := getBeginAndEnd(numOfProbes, parts, which)
	fmt.Printf("下标:%d->%d\n", begin, end)

	for i := begin; i <= end; i++ {
		wg.Add(1)
		fmt.Println("准备放入队列:", netlocsToProbe[i])
		_ = p.Invoke(string(netlocsToProbe[i]))
	}

	wg.Wait()
	fmt.Println("TCP层协议识别任务完成")
	return result
}
