// 检测给出的IP的存活性
package scan

import (
	"fmt"
	"math/rand"
	"myscanner/core/types"
	"myscanner/lib/slog"
	"myscanner/settings"
	"net"
	"sync"

	"github.com/c-robinson/iplib"
	"github.com/panjf2000/ants/v2"
)

//检查主机存活性，返回存活的主机列表，分布式模式，
//TODO: 好像没扫完就把结果返回了，这样不行
func ScanTargetsWithShuffle(t types.Targets, parts int, which int, randID int64) []string {

	aliveHosts := []string{}
	hostsToScan := []string{}
	var lock sync.Mutex
	var wg sync.WaitGroup

	threads := settings.HOST_SCAN_THREADS

	p, _ := ants.NewPoolWithFunc(threads, func(i interface{}) {
		ip := i.(string)
		if checkAlive(ip) {
			lock.Lock()
			aliveHosts = append(aliveHosts, ip)
			lock.Unlock()
		}
		wg.Done()
	})
	defer p.Release()

	//1） 把单独列出的IP加入待扫描队列
	hostsToScan = append(hostsToScan, t.IPAddrs...)

	//2） 把成段给出的IP加入待扫描队列
	for _, rangeobj := range t.IPRanges {
		start := net.ParseIP(rangeobj.Start)
		end := net.ParseIP(rangeobj.End)
		//对于每个范围，遍历该范围内的所有ip
		//ip=start; ip<=end ; ip++
		for ipobj := start; iplib.CompareIPs(ipobj, end) == -1 || iplib.CompareIPs(ipobj, end) == 0; ipobj = iplib.NextIP(ipobj) {
			// TODO: 这里需要校验ip的合法性，例如广播IP、组播IP，暂且假设每个IP都是合法的
			ip := ipobj.String()
			//slog.Debug("将" + ip + "加入存活检测队列")

			hostsToScan = append(hostsToScan, ip)
			//p.In <- ip
		}
	}

	//3) 打乱并加入队列
	numOfHosts := len(hostsToScan)

	fmt.Println("要扫描 ", numOfHosts, "个主机")

	rand.Seed(randID)
	rand.Shuffle(numOfHosts, func(i, j int) { hostsToScan[i], hostsToScan[j] = hostsToScan[j], hostsToScan[i] })

	//fmt.Println("打乱后:", hostsToScan)

	begin, end := getBeginAndEnd(numOfHosts, parts, which)

	fmt.Printf("下标:%d->%d\n", begin, end)

	for i := begin; i <= end; i++ {
		//fmt.Println(hostsToScan[i])
		wg.Add(1)
		_ = p.Invoke(string(hostsToScan[i]))
	}

	wg.Wait()
	slog.Warning("主机存活性探测任务完成")

	return aliveHosts
}
