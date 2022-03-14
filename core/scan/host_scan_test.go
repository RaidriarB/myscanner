package scan

import (
	"fmt"
	"testing"
	"time"

	"myscanner/core/types"
)

func TestLoadTargets(t *testing.T) {
	// TODO: 编写测试

	fmt.Printf("\n%#v\n\n", LoadTargets())
}

func TestScanTargets(t *testing.T) {

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

	targets2 := types.Targets{
		IPAddrs: []string{},
		IPRanges: []types.IPRange{
			{Start: "10.5.0.100", End: "10.5.0.149"}, //8个
		},
	}

	targets3 := types.Targets{
		IPAddrs: []string{
			"193.170.192.215", "150.254.36.120", "5.144.173.191",
		},
		IPRanges: []types.IPRange{},
	}

	_ = targets1
	_ = targets3

	//var randid = time.Now().UnixNano()
	var randid_s int64 = 12345

	// 别写反了
	var alivehosts []string = ScanTargetsWithShuffle(targets2, 1, 1, randid_s)

	fmt.Println("存活的主机和如下：")
	time.Sleep(2 * time.Second)
	fmt.Println(alivehosts)
	fmt.Println("存活数", len(alivehosts))

}
