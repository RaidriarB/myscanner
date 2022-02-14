package host_scan

import (
	"fmt"
	"testing"

	"myscanner/core/types"
)

func TestLoadTargets(t *testing.T) {
	// TODO: 编写测试

	fmt.Printf("\n%#v\n\n", LoadTargets(types.Config{}))
}

func TestScanTargets(t *testing.T) {
	targets := types.Targets{
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

	var alivehosts []string = ScanTargets(targets)
	fmt.Println("存活的主机如下：")
	fmt.Println(alivehosts)

}
