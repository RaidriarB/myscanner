package newnmap

import (
	"myscanner/lib/ping"
)

func HostDiscovery(ip string) bool {
	if ping.Check(ip) {
		return true
	}
	return false
}
