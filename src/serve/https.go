package serve

import (
	"crypto/tls"
	"net/http"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
	"time"
)

func ServeHTTPS() {
	tcpListener := ts.ListenTailnet(443)
	httpsListener := tls.NewListener(tcpListener, &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			start := time.Now()
			defer logger.Log("GetCertificate took %v", time.Since(start))
			return ts.GetCertificate(hello)
		},
	})
	defer httpsListener.Close()
	if err := http.Serve(httpsListener, makeProxy(false)); err != nil {
		logger.Fatal("error serving HTTPS: %v", err)
	}
}
