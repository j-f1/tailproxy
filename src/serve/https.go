package serve

import (
	"crypto/tls"
	"net/http"
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
	"time"
)

func ServeHTTPS(opts config.Options, handler http.Handler) {
	tcpListener := ts.Listen("tcp", ":443")
	httpsListener := tls.NewListener(tcpListener, &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			start := time.Now()
			defer logger.Log("GetCertificate took %v\n", time.Since(start))
			return ts.GetCertificate(hello)
		},
	})
	defer httpsListener.Close()
	if err := http.Serve(httpsListener, handler); err != nil {
		logger.Fatal("error serving HTTPS: %v\n", err)
	}
}
