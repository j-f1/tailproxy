package serve

import (
	"net/http"
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/ts"
)

func ServeHTTP(opts config.Options, handler http.Handler) {
	httpListener := ts.Listen("tcp", ":80")
	defer httpListener.Close()
	if opts.HTTPSMode == config.HTTPSRedirect {
		if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			suffix, err := ts.MagicDNSSuffix()
			if err != "" {
				http.Error(w, err, http.StatusInternalServerError)
			}
			fqdn := opts.MachineName + "." + suffix
			http.Redirect(w, r, "https://"+fqdn+r.RequestURI, http.StatusMovedPermanently)
		})); err != nil {
			logger.Fatal("error serving HTTP redirect: %v\n", err)
		}
	} else {
		if err := http.Serve(httpListener, handler); err != nil {
			logger.Fatal("error serving HTTP: %v\n", err)
		}
	}
}
