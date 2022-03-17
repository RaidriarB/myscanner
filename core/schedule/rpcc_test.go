package schedule

import (
	"fmt"
	"myscanner/core/output"
	"myscanner/core/types"
	"testing"
)

var targets = types.Targets{
	IPAddrs: []string{},
	IPRanges: []types.IPRange{
		{Start: "10.5.0.101", End: "10.5.0.110"}, //3个
	},
}

func TestScanTargetsDistributed(t *testing.T) {

	alivehosts := ScanTargetsDistributed(targets)
	fmt.Println(alivehosts)

}

func TestScanPortsDistributed(t *testing.T) {

	fmt.Println("先扫描主机")
	alivehosts := ScanTargetsDistributed(targets)
	fmt.Println("扫描主机完毕。下面扫描端口")
	aliveHostsAndPorts := ScanPortsDistributed(alivehosts)

	fmt.Println("---输出端口存活结果如下:")
	for k, v := range aliveHostsAndPorts {
		fmt.Printf("%v:%v(len:%d) \n", k, v, len(v))
	}
	fmt.Println("---输出完毕.")
}

func TestServiceProbeDistributed(t *testing.T) {

	fmt.Println("先扫描主机")
	alivehosts := ScanTargetsDistributed(targets)
	fmt.Println("扫描主机完毕。下面扫描端口")
	aliveHostsAndPorts := ScanPortsDistributed(alivehosts)
	fmt.Println("扫描端口完毕，下面服务识别")
	result := ServiceProbeDistributed(aliveHostsAndPorts)

	aliveNum := len(alivehosts)
	aliveService := len(result)

	fmt.Printf("存活主机有%d个，存活的服务有%d个\n", aliveNum, aliveService)

	output.ProcessResult(result)
}
