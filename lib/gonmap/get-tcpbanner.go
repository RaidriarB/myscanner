package gonmap

import (
	"context"
	"myscanner/lib/slog"
	"myscanner/lib/urlparse"
	"time"
)

func GetTcpBanner(netloc string, nmap *Nmap, timeout time.Duration) *TcpBanner {
	//slog.Debug(netloc)
	parse, err := urlparse.Load(netloc)
	if err != nil {
		slog.Debug(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resChan := make(chan *TcpBanner)
	go func() {
		var r TcpBanner
		defer func() {
			if err := recover(); err != nil {
				if &r != nil {
					if r.Response.Length() > 0 {
						slog.Debug(err, ",", parse.UnParse(), ",", r.Status, ",response length is :", r.Response.Length())
					}
				}
			}
		}()
		r = nmap.Scan(parse.Netloc, parse.Port)
		resChan <- &r
	}()

	for {
		select {
		case <-ctx.Done():
			close(resChan)
			return nil
		case res := <-resChan:
			close(resChan)
			return res
		}
	}
}
