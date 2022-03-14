package types

// 描述待扫描的host+port的数据结构
type TargetWithPorts map[string]([]string)

// 需要并发安全
// TODO: 大修补
//type TargetWithPorts *smap.SMap
