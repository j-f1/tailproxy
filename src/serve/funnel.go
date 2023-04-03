package serve

import (
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeFunnel() {
	httpListener := ts.ListenFunnel(443)
	defer httpListener.Close()
	if err := http.Serve(httpListener, makeProxy(true)); err != nil {
		logger.Fatal("error serving to Funnel: %v", err)
	}
}
