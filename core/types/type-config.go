// 存储所有配置信息，该结构不同于Args，后者仅是前者的子集
package types

// TODO: 需要不停补充
type Config struct {
	IPRangeGiven  []IPRange //直接给出的扫描范围
	IPGivenInFile []string  //文件中给出的IP

	IPAggrSegB      bool // 是否聚合B段
	IPAggrSegC      bool // 是否聚合C段，都开启按照B段来聚合
	IPAggrThreshold int  //进行聚合的阈值（高于该阈值，跳过该网段的扫描）
}