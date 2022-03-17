package schedule

import (
	"github.com/smallnest/rpcx/server"
)

func StartServer(serverAddr string) {
	s := server.NewServer()
	s.RegisterName("ScanUnit", new(ScanUnit), "")
	err := s.Serve("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

}
