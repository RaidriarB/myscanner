// 检测给出的IP的存活性
package host_scan

import (
	"myscanner/core/types"
	"myscanner/lib/pool"
	"myscanner/lib/slog"
	"net"
	"strings"
	"sync"

	"github.com/c-robinson/iplib"
)

//预处理需要扫描存活性的主机
func LoadTargets(c types.Config) types.Targets {

	//TODO: 读取信息，进行聚合操作，生成最终要扫描的

	return types.Targets{
		IPAddrs: []string{
			"10.2.1.3", "10.2.1.4", "10.2.1.6", "10.2.1.7", "10.2.1.10", "10.2.1.13",
			"10.2.3.4", "10.2.3.5", "10.2.3.8", "10.2.3.9",
			"10.2.4.12", "10.2.4.15",
		},
		IPRanges: []types.IPRange{
			{Start: "10.5.0.0", End: "10.5.255.255"},
			{Start: "10.6.7.0", End: "10.6.7.255"},
		},
	}

}

//检查主机存活性，返回存活的主机列表
func ScanTargets(t types.Targets) []string {
	// FIXME: 需不需要Config参数？
	aliveHosts := []string{}
	// TODO: 把这个10改成配置
	var p = pool.NewPool(10)

	//主机存活性检测——设置pool中要执行的函数
	p.Function = func(i interface{}) interface{} {
		ip := i.(string)
		if checkAlive(ip) {
			return ip
		}
		return nil
	}

	//主机存活性探测——输出调度
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

	//主机存活性探测——将要检测的IP加入队列
	go func() {
		//1. 把单独列出的IP加入队列
		for _, ip := range t.IPAddrs {
			slog.Debug("将" + ip + "加入存活检测队列")
			p.In <- ip
		}
		//2. 把成段给出的IP加入队列
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

	//开始执行主机存活性探测任务
	p.Run()
	slog.Warning("主机存活性探测任务完成")
	return aliveHosts
}

func checkAlive(ip string) bool {

	//FIXME: 目前只是模拟一下

	//fmt.Println("检查IP[" + ip + "]的存活性")

	if strings.HasSuffix(ip, "2") || strings.HasSuffix(ip, "4") || strings.HasSuffix(ip, "6") || strings.HasSuffix(ip, "8") || strings.HasSuffix(ip, "0") {
		return true
	}
	//return gonmap.HostDiscovery(ip)
	return false
}
