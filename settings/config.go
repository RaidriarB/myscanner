package settings

import "myscanner/core/types"

const (
	DEV_MODE = true // 开发模式
	MAXPORT  = 16   // 端口从1扫到多少，正常设为65536
)

const (

	//TODO: implement
	SCAN_STEPS = 3 // 1:仅检测主机存活性 2:仅探测主机和端口 3:并且识别banner

	SCANALL = true // 是否扫描所有端口

	CAREFUL_MODE = true //TODO: implement 不扫描特别常见的敏感端口

	PORTLIST_MODE = 2 // 1: 空端口列表 2: debug的简单列表 3:建议列表 4:top1000列表 5:完整列表

	PORT_SCAN_THREADS     = 6
	HOST_SCAN_THREADS     = 2
	SERVICE_PROBE_THREADS = 10
	PORT_SCAN_TIMEOUT     = 2 * 1000000000 // nanoseconds

	KNOWN_IP_FILE = "settings/known_IP.txt"

	SAVE_RESULT_TO_FILE = true
	SAVE_RESULT_TO_DB   = false
	RESULT_FILE         = "output.txt"

	// TODO: implement
	RESULT_DB_FILE = "output.sqlite"
)

var IP_RANGES_TO_SCAN = []types.IPRange{
	//{Start: "10.3.244.230", End: "10.3.244.255"},
}
