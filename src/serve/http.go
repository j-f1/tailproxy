package serve

import (
	"context"
	"net"
	"net/http"
	"os"
	"tailproxy/src/config"
	"tailproxy/src/logger"

	"tailscale.com/tsnet"
)

func ServeHTTP(s *tsnet.Server, opts config.Options, handler http.Handler) {
	httpListener, err := s.Listen("tcp", ":80")
	if err != nil {
		logger.Err("error listening on port 80: %v\n", err)
		os.Exit(1)
	}
	genericServeHTTP(httpListener, opts, handler)
}

func genericServeHTTP(httpListener net.Listener, opts config.Options, handler http.Handler) {
	defer httpListener.Close()
	if opts.HTTPSMode == config.HTTPSRedirect {
		if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status, err := lc.Status(context.Background())
			if err != nil || status == nil {
				logger.Err("error getting profile status: %v\n", err)
				http.Error(w, "error getting profile status", http.StatusInternalServerError)
				return
			}
			if status.CurrentTailnet == nil {
				logger.Err("not logged in (CurrentTailnet is nil)\n")
				http.Error(w, "not logged in (CurrentTailnet is nil)", http.StatusForbidden)
				return
			}
			fqdn := opts.MachineName + "." + status.CurrentTailnet.MagicDNSSuffix
			http.Redirect(w, r, "https://"+fqdn+r.RequestURI, http.StatusMovedPermanently)
		})); err != nil {
			logger.Err("error serving HTTP redirect: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := http.Serve(httpListener, handler); err != nil {
			logger.Err("error serving HTTP: %v\n", err)
			os.Exit(1)
		}
	}
}
