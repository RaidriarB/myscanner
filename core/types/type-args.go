// 存储命令行参数的数据结构
package types

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	//TODO: implement
	TargetRanges string // e.g. "10.0.0.1-10.0.0.20,10.0.3.0-10.0.5.255"
	PortRanges   string // e.g. "1-100,50000-60000" 默认为1-65535
}

//初始化参数
func (a *Args) init() {
	a.define()
	flag.Parse()
	a.checkArgs()
}

//定义参数
func (a *Args) define() {
	flag.StringVar(&a.TargetRanges, "t", "", `待扫描的目标，可以指定多个IP范围。例如："10.0.0.1-10.0.0.20;10.0.3.0-10.0.5.255"`)
	flag.StringVar(&a.PortRanges, "p", "1-65535", `待扫描的端口范围，可以指定多个端口范围。例如："1-100,50000-60000"。默认为1-65535`)
}

//检查命令行参数存在与否、是否冲突意义的合法性。暂时不检查语义上的合法性。
func (a *Args) checkArgs() {
	if a.TargetRanges == "" {
		fmt.Println("待扫描的目标不能为空！")
		os.Exit(0)
	}
}
