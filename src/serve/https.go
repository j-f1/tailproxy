package serve

import (
	"crypto/tls"
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeHTTPS() {
	tcpListener := ts.ListenTailnet(443)
	tlsListener := tls.NewListener(tcpListener, &tls.Config{
		GetCertificate: ts.GetCertificate,
	})
	defer tlsListener.Close()
	if err := http.Serve(tlsListener, makeProxy(false)); err != nil {
		logger.Fatal("error serving HTTPS: %v", err)
	}
}
