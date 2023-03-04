package serve

import (
	"net/http"
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	suffix, err := ts.MagicDNSSuffix(r.Context())
	if err != "" {
		http.Error(w, err, http.StatusInternalServerError)
	}
	fqdn := config.MachineName + "." + suffix
	http.Redirect(w, r, "https://"+fqdn+r.RequestURI, http.StatusFound)
}

func ServeHTTP() {
	httpListener := ts.Listen("tcp", ":80")
	defer httpListener.Close()
	if config.HTTPSMode == config.HTTPSRedirect {
		if err := http.Serve(httpListener, http.HandlerFunc(redirectToHTTPS)); err != nil {
			logger.Fatal("error serving HTTP redirect: %v\n", err)
		}
	} else {
		if err := http.Serve(httpListener, makeProxy()); err != nil {
			logger.Fatal("error serving HTTP: %v\n", err)
		}
	}
}
