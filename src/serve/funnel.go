package serve

import (
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeFunnel() {
	tcpListener := ts.ListenFunnel(443)
	defer tcpListener.Close()
	if err := http.Serve(tcpListener, makeProxy(true)); err != nil {
		logger.Fatal("error serving to Funnel: %v", err)
	}
}
