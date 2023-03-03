package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/serve"
	"time"

	"tailscale.com/tsnet"
)

func main() {
	opts := config.ParseOptions()

	logger.Log("machine name: %v, target: %v, https mode: %v\n", opts.MachineName, opts.Target, opts.HTTPSMode)

	s := new(tsnet.Server)
	s.Hostname = opts.MachineName
	s.Ephemeral = true

	if err := s.Start(); err != nil {
		logger.Err("error starting server: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	lc, err := s.LocalClient()
	if err != nil {
		logger.Err("error getting local client: %v\n", err)
		os.Exit(1)
	}

	if lc == nil {
		logger.Err("no local client; are you running tailscaled?\n")
		os.Exit(1)
	}

	err = lc.StartLoginInteractive(context.Background())
	if err != nil {
		logger.Err("error starting login: %v\n", err)
		os.Exit(1)
	}
	status, err := lc.Status(context.Background())
	if err != nil || status == nil {
		logger.Err("error getting profile status: %v\n", err)
		os.Exit(1)
	}

	proxy := serve.MakeProxy(opts)

	go func() {
	}()

	if opts.HTTPSMode != config.HTTPSOff {
		var httpsListener net.Listener
		tcpListener, err := s.Listen("tcp", ":443")
		if err != nil {
			logger.Err("error listening on port 443: %v\n", err)
			os.Exit(1)
		}
		httpsListener = tls.NewListener(tcpListener, &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				start := time.Now()
				defer logger.Log("GetCertificate took %v\n", time.Since(start))
				return lc.GetCertificate(hello)
			},
		})
		go func() {
			defer httpsListener.Close()
			if err := http.Serve(httpsListener, proxy); err != nil {
				logger.Err("error serving HTTPS: %v\n", err)
				os.Exit(1)
			}
		}()
	}

	logger.Log("listening as %s, forwarding to %v\n", opts.MachineName, opts.Target)

	select {}
}
