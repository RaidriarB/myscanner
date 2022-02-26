// 检测给出的IP的存活性
package host_scan

import (
	"bufio"
	"fmt"
	"io"
	"myscanner/core/types"
	"myscanner/lib/gonmap"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
	"myscanner/settings"
	"net"
	"os"
	"sync"

	"github.com/c-robinson/iplib"
)

//预处理需要扫描存活性的主机
func LoadTargets(c types.Config) types.Targets {

	//TODO: 读取信息，进行聚合操作，生成最终要扫描的
	//targets := types.Targets{}
	//TODO: 改成配置，不要硬编码

	ipaddrs := []string{}
	ipListFile, err := os.Open("../../settings/known_IP.txt")
	defer ipListFile.Close()
	if err != nil {
		fmt.Printf("known_IP.txt文件不存在")
	}

	br := bufio.NewReader(ipListFile)
	lineNumOfFile := 0
	for {
		lineNumOfFile++
		ip, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if net.ParseIP(string(ip)) == nil {
			fmt.Printf("IP列表格式有误：第%d行%s，不是合法IP格式", lineNumOfFile, ip)

		} else {
			ipaddrs = append(ipaddrs, string(ip))
		}
	}

	//TODO: 读取范围IP
	ipranges := []types.IPRange{
		{Start: "10.5.0.0", End: "10.5.255.255"},
		{Start: "10.6.7.0", End: "10.6.7.255"},
	}

	return types.Targets{
		IPAddrs:  ipaddrs,
		IPRanges: ipranges,
	}
}

//检查主机存活性，返回存活的主机列表
func ScanTargets(t types.Targets) []string {
	// FIXME: 需不需要Config作为参数传入？
	aliveHosts := []string{}
	var p = pool.NewPool(settings.HOST_SCAN_THREADS)
	//1.设置pool中要执行的函数
	p.Function = func(i interface{}) interface{} {
		ip := i.(string)
		if checkAlive(ip) {
			return ip
		}
		return nil
	}

	//2. 输出调度
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			for out := range p.Out {
				if out != nil {
					ip := (out).(string)
					aliveHosts = append(aliveHosts, ip)
				}
			}
			wg.Done()
		}()
		wg.Wait()
	}()

	//3. 将要检测的IP加入队列
	go func() {
		//1） 把单独列出的IP加入队列
		for _, ip := range t.IPAddrs {
			slog.Debug("将" + ip + "加入存活检测队列")
			p.In <- ip
		}
		//2） 把成段给出的IP加入队列
		for _, rangeobj := range t.IPRanges {
			start := net.ParseIP(rangeobj.Start)
			end := net.ParseIP(rangeobj.End)
			//对于每个范围，遍历该范围内的所有ip
			for ipobj := start; iplib.CompareIPs(ipobj, end) == -1; ipobj = iplib.NextIP(ipobj) {
				// TODO: 这里需要校验ip的合法性，例如广播IP、组播IP，暂且假设每个IP都是合法的
				ip := ipobj.String()
				slog.Debug("将" + ip + "加入存活检测队列")
				p.In <- ip
			}
		}
		//关闭主机存活性探测下发信道
		slog.Info("主机存活性探测任务下发完毕")
		p.InDone()
	}()

	//4.开始执行主机存活性探测任务
	p.Run()
	slog.Warning("主机存活性探测任务完成")
	return aliveHosts
}

func checkAlive(ip string) bool {

	return gonmap.HostDiscovery(ip)
	//return ping.Check(ip)
}
