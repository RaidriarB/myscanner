package schedule

import (
	"fmt"
	"myscanner/settings"
	"sync"
	"testing"
)

func TestStartServer(t *testing.T) {

	// 创建很多个服务

	var wg sync.WaitGroup

	for i, addr := range settings.SRVADDRS {
		wg.Add(1)
		go func(i int, addr string) {
			StartServer(addr)
			fmt.Printf("第%d个服务建立成功。\n", i)
			wg.Done()
		}(i, addr)
	}
	wg.Wait()

}
