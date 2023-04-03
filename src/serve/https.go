package serve

import (
	"crypto/tls"
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeHTTPS() {
	tcpListener := ts.ListenTailnet(443)
	httpsListener := tls.NewListener(tcpListener, &tls.Config{
		GetCertificate: ts.GetCertificate,
	})
	defer httpsListener.Close()
	if err := http.Serve(httpsListener, makeProxy(false)); err != nil {
		logger.Fatal("error serving HTTPS: %v", err)
	}
}
