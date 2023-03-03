package serve

import (
	"net/http"
	_ "net/http/pprof"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServePProf() {
	httpListener := ts.Listen("tcp", ":6060")
	defer httpListener.Close()
	if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/debug/pprof/", http.StatusFound)
		} else {
			http.DefaultServeMux.ServeHTTP(w, r)
		}
	})); err != nil {
		logger.Fatal("http.Serve: %v", err)
	}
}
