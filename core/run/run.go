package run

import (
	"fmt"
	"myscanner/core/output"
	"myscanner/core/scan"
	"time"
)

func StartScan() {
	targets := scan.LoadTargets()

	//1. 主机存活检测
	fmt.Println("正在进行主机存活性检测...")
	aliveHosts := scan.ScanTargets(targets)
	time.Sleep(time.Microsecond * 200)
	aliveNum := len(aliveHosts)

	//2. 端口存活检测
	fmt.Println("正在进行端口存活性检测...")
	aliveHostsAndPorts := scan.ScanPorts(aliveHosts)
	time.Sleep(time.Microsecond * 200)

	//3. 服务识别
	fmt.Println("正在进行服务识别...")
	result := scan.ServiceProbe(aliveHostsAndPorts)
	time.Sleep(time.Microsecond * 200)
	aliveService := len(result)

	fmt.Printf("存活主机有%d个，存活的服务有%d个\n", aliveNum, aliveService)

	output.ProcessResult(result)

}
