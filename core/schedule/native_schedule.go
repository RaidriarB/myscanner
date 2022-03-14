package schedule

import (
	"fmt"
	"myscanner/core/output"
	"myscanner/core/scan"
	"myscanner/core/types"
	"myscanner/settings"
	"sync"
	"time"
)

//在本地模拟分布式
func NativeScan(parts int) {

	var targets = scan.LoadTargets()
	var randID int64 = 1234567

	//1. 主机存活检测
	fmt.Println("正在进行主机存活性检测...")
	var aliveHosts []string
	var wg1 sync.WaitGroup
	var lock1 sync.Mutex
	//这个循环写的对吗？
	for which := 1; which <= parts; which++ {
		fmt.Printf("正在安排第%d/%d个部分\n", which, parts)
		wg1.Add(1)
		go func(w int) {
			aliveHosts_part := scan.ScanTargetsWithShuffle(targets, parts, w, randID)
			lock1.Lock()
			aliveHosts = append(aliveHosts, aliveHosts_part...)
			lock1.Unlock()

			fmt.Printf("第%d个部分扫描完成。\n", w)
			wg1.Done()
		}(which)
	}
	wg1.Wait()
	time.Sleep(time.Microsecond * 200)
	aliveNum := len(aliveHosts)

	fmt.Printf("存活主机:%d", aliveNum)
	fmt.Println(aliveHosts)

	//2. 端口存活检测
	fmt.Println("正在进行端口存活性检测...")
	var aliveHostsAndPorts sync.Map
	var wg2 sync.WaitGroup
	var lock2 sync.Mutex

	var portlist []int
	var scanall bool
	if settings.DEV_MODE {
		portlist = settings.PORTLIST_FOR_DEBUG
		scanall = true
	} else {
		portlist = settings.PORTLIST_SIMPLIFIED
		scanall = false
	}

	//这个循环写的对吗？
	for which := 1; which <= parts; which++ {
		fmt.Printf("正在安排第%d/%d个部分\n", which, parts)
		wg2.Add(1)
		go func(w int) {

			aliveHostsAndPorts_part := scan.ScanPortsWithShuffle(aliveHosts, portlist, scanall, parts, w, randID)

			//将部分整合
			lock2.Lock()
			aliveHostsAndPorts_part.Range(func(k, v interface{}) bool {
				lst := []string{}
				if v, ok := aliveHostsAndPorts.Load(k); ok {
					lst = v.([]string)
				}
				lst = append(lst, v.([]string)...)
				aliveHostsAndPorts.Store(k, lst)
				return true
			})
			lock2.Unlock()

			fmt.Printf("第%d个部分扫描完成。\n", w)
			wg2.Done()
		}(which)
	}
	wg2.Wait()

	fmt.Println("---输出端口存活结果如下:")
	aliveHostsAndPorts.Range(func(k, v interface{}) bool {
		fmt.Printf("%v:%v(len:%d) \n", k, v, len(v.([]string)))
		return true
	})
	fmt.Println("---输出完毕.")

	time.Sleep(time.Microsecond * 500)

	//3. 服务识别
	fmt.Println("正在进行服务识别...")
	var result types.TargetPortBanners
	var wg3 sync.WaitGroup
	var lock3 sync.Mutex
	//这个循环写的对吗？
	scan.InitNmap()

	for which := 1; which <= parts; which++ {
		wg3.Add(1)
		go func(w int) {

			// FIXME: 这里好像有并发问题
			result_part := scan.ServiceProbeWithShuffle(aliveHostsAndPorts, parts, w, randID)

			lock3.Lock()
			result = append(result, result_part...)
			lock3.Unlock()

			wg3.Done()
		}(which)
	}
	wg3.Wait()
	time.Sleep(time.Microsecond * 200)

	aliveService := len(result)

	fmt.Printf("存活主机有%d个，存活的服务有%d个\n", aliveNum, aliveService)

	output.ProcessResult(result)

}
