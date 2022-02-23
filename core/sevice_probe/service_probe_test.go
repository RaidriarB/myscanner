package service_probe

import (
	"fmt"
	"myscanner/core/types"
	"testing"
)

func TestServiceProbe(t *testing.T) {
	twp := types.TargetWithPorts{
		"193.170.192.215": {"8080", "80"},
		"150.254.36.120":  {"22", "80", "8080"},
		"5.144.173.191":   {"21", "22", "80", "443", "3306", "8009", "8080"},
	}

	twp2 := types.TargetWithPorts{
		"5.144.173.191": {"8009", "8080", "22", "80"},
	}
	_ = twp2
	//_ = twp2

	// targets := types.Targets{
	// 	IPAddrs: []string{
	// 		"193.170.192.215", "150.254.36.120", "5.144.173.191",
	// 	},
	// 	IPRanges: []types.IPRange{},
	// }
	//var alivehosts []string = host_scan.ScanTargets(targets)
	//var alivehostandports = port_scan.ScanPorts(alivehosts)
	//result := ServiceProbe(alivehostandports)

	result := ServiceProbe(twp)
	for _, r := range result {
		fmt.Println(r.Target)
	}

}
