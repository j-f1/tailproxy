package serve

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"tailproxy/src/config"
)

func MakeProxy(opts config.Options) http.Handler {
	var start time.Time
	return &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			start = time.Now()
			fmt.Printf("tailproxy: %v %v %v\n", r.In.RemoteAddr, r.In.Method, r.In.URL)
			r.SetXForwarded()
			r.SetURL(opts.Target)
			r.Out.Host = r.In.Host
		},
		ModifyResponse: func(r *http.Response) error {
			fmt.Printf("tailproxy: %v %v %v %v %v\n", r.Request.RemoteAddr, r.Request.Method, r.Request.URL, r.StatusCode, time.Since(start))
			return nil
		},
	}
}
