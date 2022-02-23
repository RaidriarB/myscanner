package port_scan

import (
	"fmt"
	"myscanner/core/host_scan"
	"myscanner/core/types"
	"testing"
)

func TestScanPorts(t *testing.T) {

	targets1 := types.Targets{
		IPAddrs: []string{
			"10.2.1.3", "10.2.1.4", "10.2.1.6", "10.2.1.7", "10.2.1.10", "10.2.1.13",
			"10.2.3.4", "10.2.3.5", "10.2.3.8", "10.2.3.9",
			"10.2.4.12", "10.2.4.15",
		},
		IPRanges: []types.IPRange{
			{Start: "10.5.254.252", End: "10.5.255.6"},
			{Start: "10.6.7.0", End: "10.6.7.7"},
		},
	}
	_ = targets1

	targets := types.Targets{
		IPAddrs: []string{
			"193.170.192.215", "150.254.36.120", "5.144.173.191",
		},
		IPRanges: []types.IPRange{},
	}

	var alivehosts []string = host_scan.ScanTargets(targets)
	var result = ScanPorts(alivehosts)
	fmt.Println("存活的主机和端口如下：")
	fmt.Println(result)

}
