package schedule

import (
	"myscanner/core/types"
	"testing"
)

func TestNativeSchedule(t *testing.T) {
	var targets = types.Targets{
		IPAddrs: []string{},
		IPRanges: []types.IPRange{
			{Start: "10.5.0.101", End: "10.5.0.110"}, //10个
		},
	}

	NativeScan(6, targets)

}
