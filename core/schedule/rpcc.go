package schedule

import (
	"context"
	"fmt"
	"myscanner/core/types"
	"myscanner/settings"
	"sync"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

//分布式
func ScanTargetsDistributed(targets types.Targets) []string {
	//var targets = scan.LoadTargets()
	var randID int64 = 1234567
	var parts = len(settings.SRVADDRS)

	//1. 主机存活检测
	fmt.Println("正在进行主机存活性检测...")
	var aliveHosts []string
	var wg1 sync.WaitGroup
	var lock1 sync.Mutex

	// TODO: 找出哪几个服务是开着的
	var aliveServiceAddrs = settings.SRVADDRS //暂时先用这个额
	parts = len(aliveServiceAddrs)

	for i, addr := range aliveServiceAddrs {
		var which = i + 1 // 从1开始
		fmt.Printf("正在安排第%d/%d个部分\n", which, parts)
		wg1.Add(1)
		go func(w int, addr string) {

			//rpc-参数与调用
			args := ScanTargetsArgs{
				Targets: targets,
				Parts:   parts,
				Which:   w,
				RandID:  randID,
			}
			fmt.Printf("args的值: %#v\n", args)

			reply := &ScanTargetsReply{}
			CallRPC(addr, "ScanUnit", "ScanTargets", args, reply)

			fmt.Printf("第%d个的回复:%#v\n", w, reply)

			aliveHosts_part := reply.AliveHosts
			// 将一个part的结果合并
			lock1.Lock()
			aliveHosts = append(aliveHosts, aliveHosts_part...)
			lock1.Unlock()

			fmt.Printf("第%d个部分扫描完成。\n", w)
			wg1.Done()
		}(which, addr)
	}
	wg1.Wait()
	time.Sleep(time.Microsecond * 200)
	aliveNum := len(aliveHosts)

	fmt.Printf("存活主机:%d", aliveNum)
	fmt.Println(aliveHosts)

	return aliveHosts

}

func ScanPortsDistributed(aliveHosts []string) types.TargetWithPorts {
	var randID int64 = 1234567
	var parts = len(settings.SRVADDRS)

	//2. 端口存活检测
	fmt.Println("正在进行端口存活性检测...")
	var aliveHostsAndPorts = make(types.TargetWithPorts)
	var wg2 sync.WaitGroup
	var lock2 sync.Mutex

	var portlist []int
	var scanall bool
	if settings.DEV_MODE {
		portlist = settings.PORTLIST_FOR_DEBUG
		scanall = true
	} else {
		portlist = settings.PORTLIST_SIMPLIFIED
		scanall = false
	}

	var aliveServiceAddrs = settings.SRVADDRS //暂时先用这个额
	parts = len(aliveServiceAddrs)
	for i, addr := range aliveServiceAddrs {

		var which = i + 1 // 从1开始
		fmt.Printf("正在安排第%d/%d个部分\n", which, parts)
		wg2.Add(1)
		go func(w int, addr string) {

			//rpc-参数与调用
			args := ScanPortsArgs{
				AliveHosts: aliveHosts,
				Portlist:   portlist,
				SCANALL:    scanall,
				Parts:      parts,
				Which:      w,
				RandID:     randID,
			}
			fmt.Printf("args的值: %#v\n", args)

			reply := &ScanPortsReply{}
			CallRPC(addr, "ScanUnit", "ScanPorts", args, reply)

			fmt.Printf("第%d个的回复:%#v\n", w, reply)

			aliveHostsAndPorts_part := reply.UpHostsWithPorts

			// 将一个part的结果合并
			lock2.Lock()
			for k, v := range aliveHostsAndPorts_part {
				aliveHostsAndPorts[k] = append(aliveHostsAndPorts[k], v...)
			}
			lock2.Unlock()

			fmt.Printf("第%d个部分扫描完成。\n", w)
			wg2.Done()
		}(which, addr)
	}
	wg2.Wait()

	fmt.Println("---输出端口存活结果如下:")
	for k, v := range aliveHostsAndPorts {
		fmt.Printf("%v:%v(len:%d) \n", k, v, len(v))
	}
	fmt.Println("---输出完毕.")
	time.Sleep(time.Microsecond * 500)
	return aliveHostsAndPorts

}

func ServiceProbeDistributed(aliveHostsAndPorts types.TargetWithPorts) types.TargetPortBanners {

	var randID int64 = 1234567
	var parts = len(settings.SRVADDRS)
	//3. 服务识别
	fmt.Println("正在进行服务识别...")
	var result types.TargetPortBanners
	var wg3 sync.WaitGroup
	var lock3 sync.Mutex
	//这个循环写的对吗？
	//TODO: 这个初始化必须放在Probe时，这样不行！

	// TODO: 找出哪几个服务是开着的
	var aliveServiceAddrs = settings.SRVADDRS //暂时先用这个额
	parts = len(aliveServiceAddrs)

	for i, addr := range aliveServiceAddrs {
		var which = i + 1 // 从1开始
		wg3.Add(1)
		go func(w int, addr string) {

			// //rpc-准备工作
			// d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
			// opt := client.DefaultOption
			// opt.SerializeType = protocol.JSON
			// // TODO: 这里需要改个名字。
			// xclient := client.NewXClient("ScanUnit", client.Failtry, client.RandomSelect, d, opt)
			// defer xclient.Close()

			//rpc-参数与调用
			args := ServiceProbeArgs{
				Twp:    aliveHostsAndPorts,
				Parts:  parts,
				Which:  w,
				RandID: randID,
			}
			fmt.Printf("args的值: %#v\n", args)

			reply := &ServiceProbeReply{}
			CallRPC(addr, "ScanUnit", "ServiceProbe", args, reply)

			fmt.Printf("第%d个的回复:%#v\n", w, reply)

			result_part := reply.Result

			// FIXME: 这里好像有并发问题
			//result_part := scan.ServiceProbeWithShuffle(aliveHostsAndPorts, parts, w, randID)

			lock3.Lock()
			result = append(result, result_part...)
			lock3.Unlock()

			wg3.Done()
		}(which, addr)
	}
	wg3.Wait()

	return result

}

func CallRPC(addr string, path string, name string, args interface{}, reply interface{}) {
	//1. 准备工作
	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	opt := client.DefaultOption
	opt.SerializeType = protocol.MsgPack
	xclient := client.NewXClient(path, client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()

	//2. 调用
	err := xclient.Call(context.Background(), name, args, reply)
	if err != nil {
		//TODO: error handling
		fmt.Printf("failed to call: %v", err)
	}

}
