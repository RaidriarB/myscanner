// 在主机扫描中，刻画待扫描目标的数据结构，分为单个IP点和IP段。
package types

type IPRange struct {
	Start, End string
}

type Targets struct {
	IPAddrs  []string
	IPRanges []IPRange
}
