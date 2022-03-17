package schedule

import (
	"context"
	"myscanner/core/scan"
	"myscanner/core/types"
)

// 千万记住要大写！！
type ScanTargetsArgs struct {
	Targets types.Targets
	Parts   int
	Which   int
	RandID  int64
}

type ScanTargetsReply struct {
	AliveHosts []string
}

type ScanPortsArgs struct {
	AliveHosts []string
	Portlist   []int
	SCANALL    bool
	Parts      int
	Which      int
	RandID     int64
}

type ScanPortsReply struct {
	UpHostsWithPorts types.TargetWithPorts
}

type ServiceProbeArgs struct {
	Twp    types.TargetWithPorts
	Parts  int
	Which  int
	RandID int64
}

type ServiceProbeReply struct {
	Result types.TargetPortBanners
}

// TODO: 这个结构咋写比较好
type ScanUnit struct{}

func (s *ScanUnit) ScanTargets(ctx context.Context, args ScanTargetsArgs, reply *ScanTargetsReply) error {

	tar := args.Targets
	pr := args.Parts
	w := args.Which
	ra := args.RandID

	al := scan.ScanTargetsWithShuffle(tar, pr, w, ra)

	reply.AliveHosts = al

	return nil
}

func (s *ScanUnit) ScanPorts(ctx context.Context, args ScanPortsArgs, reply *ScanPortsReply) error {

	ah := args.AliveHosts
	pl := args.Portlist
	all := args.SCANALL
	pr := args.Parts
	w := args.Which
	ra := args.RandID

	result := scan.ScanPortsWithShuffle(ah, pl, all, pr, w, ra)

	reply.UpHostsWithPorts = result

	return nil
}

func (s *ScanUnit) ServiceProbe(ctx context.Context, args ServiceProbeArgs, reply *ServiceProbeReply) error {
	twp := args.Twp
	pr := args.Parts
	w := args.Which
	ra := args.RandID

	result := scan.ServiceProbeWithShuffle(twp, pr, w, ra)
	reply.Result = result
	return nil

}
