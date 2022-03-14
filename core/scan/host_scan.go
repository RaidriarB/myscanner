// 检测给出的IP的存活性
package scan

import (
	"bufio"
	"fmt"
	"io"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/slog"
	"myscanner/settings"
	"net"
	"os"
	"sync"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/panjf2000/ants/v2"
)

//预处理需要扫描存活性的主机
func LoadTargets() types.Targets {

	ipaddrs := []string{}
	ipListFile, err := os.Open(settings.KNOWN_IP_FILE)
	if err != nil {
		fmt.Printf("[%s]文件不存在,跳过此步骤\n", settings.KNOWN_IP_FILE)
	} else {
		defer ipListFile.Close()
		br := bufio.NewReader(ipListFile)
		lineNumOfFile := 0
		for {
			lineNumOfFile++
			ip, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			if net.ParseIP(string(ip)) == nil {
				fmt.Printf("IP列表格式有误：第%d行%s，不是合法IP格式\n", lineNumOfFile, ip)

			} else {
				ipaddrs = append(ipaddrs, string(ip))
			}
		}
	}
	ipranges := settings.IP_RANGES_TO_SCAN
	//TODO: 检验合法性

	return types.Targets{
		IPAddrs:  ipaddrs,
		IPRanges: ipranges,
	}
}

//检查主机存活性，返回存活的主机列表
func ScanTargets(t types.Targets) []string {
	return ScanTargetsWithShuffle(t, 1, 1, time.Now().UnixNano())
}

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
	//fmt.Println("总共要扫描 ", numOfHosts, "个主机,打乱中...")

	shuffleStringArray(hostsToScan, randID)
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

func checkAlive(ip string) bool {
	fmt.Printf("checking %s\n", ip)
	if settings.DEV_MODE {
		time.Sleep(1 * time.Second)
		return true
	} else {
		return gonmap.HostDiscovery(ip)
		//return ping.Check(ip)
	}
}
