package types

import "myscanner/lib/gonmap"

//描述扫描得到的主机+端口+端口返回的banner的数据结构
type TargetPortBanners [](*gonmap.TcpBanner)
