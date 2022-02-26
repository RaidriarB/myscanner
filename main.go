package main

import (
	"fmt"
	"myscanner/core/host_scan"
	"myscanner/core/port_scan"
	"myscanner/core/service_probe"
	"myscanner/core/types"
)

func add(a, b int) int {
	return a + b
}

func main() {
	targets := types.Targets{
		IPAddrs: []string{
			"5.144.173.191",
		},
		IPRanges: []types.IPRange{},
	}
	var alivehosts []string = host_scan.ScanTargets(targets)
	fmt.Println(alivehosts)
	var alivehostandports = port_scan.ScanPorts(alivehosts)
	fmt.Println(alivehostandports)
	result := service_probe.ServiceProbe(alivehostandports)

	//result := ServiceProbe(twp)
	for _, r := range result {
		fmt.Println(r.Target)
	}

}
