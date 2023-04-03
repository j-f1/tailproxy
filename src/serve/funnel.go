package serve

import (
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeFunnel() {
	httpListener := ts.ListenFunnel(443)
	defer httpListener.Close()
	if err := http.Serve(httpListener, makeProxy()); err != nil {
		logger.Fatal("error serving HTTP: %v\n", err)
	}
}
