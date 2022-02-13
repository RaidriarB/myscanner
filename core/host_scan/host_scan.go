// 生成待扫描的IP
package host_scan

import (
	"fmt"
	"myscanner/core/types"
	"net"

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
	// TODO: implement 现在只是一个小demo
	// FIXME: 需不需要Config参数？

	aliveHosts := []string{}

	//处理单独的IP
	for _, ip := range t.IPAddrs {
		if checkAlive(ip) {
			aliveHosts = append(aliveHosts, ip)
		}
	}

	//处理成段的IP
	for _, rangeobj := range t.IPRanges {
		start := net.ParseIP(rangeobj.Start)
		end := net.ParseIP(rangeobj.End)
		// TODO: 校验合法性。暂且假设都合法
		for ipobj := start; iplib.CompareIPs(ipobj, end) == -1; ipobj = iplib.NextIP(ipobj) {
			ip := ipobj.String()
			fmt.Println(ipobj.String())
			if checkAlive(ip) {
				aliveHosts = append(aliveHosts, ip)
			}
		}
	}
	return aliveHosts
	// return []string{"10.2.1.4", "10.2.1.6", "10.2.3.5", "10.5.0.3", "10.6.7.42"}
}

func checkAlive(ip string) bool {
	//ICMP或TCP或两者都用
	return true
	//return gonmap.HostDiscovery(ip)
}
