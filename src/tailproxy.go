package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"tailproxy/src/config"
	"tailproxy/src/serve"
	"time"

	"tailscale.com/tsnet"
)

func log(message string, args ...interface{}) {
	fmt.Printf("tailproxy: "+message, args...)
}
func logErr(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "tailproxy: "+message, args...)
}

func main() {
	opts := config.ParseOptions()

	log("machine name: %v, target: %v, https mode: %v\n", opts.MachineName, opts.Target, opts.HTTPSMode)

	s := new(tsnet.Server)
	s.Hostname = opts.MachineName
	s.Ephemeral = true

	if err := s.Start(); err != nil {
		logErr("error starting server: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	lc, err := s.LocalClient()
	if err != nil {
		logErr("error getting local client: %v\n", err)
		os.Exit(1)
	}

	if lc == nil {
		logErr("no local client; are you running tailscaled?\n")
		os.Exit(1)
	}

	err = lc.StartLoginInteractive(context.Background())
	if err != nil {
		logErr("error starting login: %v\n", err)
		os.Exit(1)
	}
	status, err := lc.Status(context.Background())
	if err != nil || status == nil {
		logErr("error getting profile status: %v\n", err)
		os.Exit(1)
	}

	proxy := serve.MakeProxy(opts)

	go func() {
		httpListener, err := s.Listen("tcp", ":80")
		if err != nil {
			logErr("error listening on port 80: %v\n", err)
			os.Exit(1)
		}

		defer httpListener.Close()
		if opts.HTTPSMode == config.HTTPSRedirect {
			if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				status, err := lc.Status(context.Background())
				if err != nil || status == nil {
					logErr("error getting profile status: %v\n", err)
					http.Error(w, "error getting profile status", http.StatusInternalServerError)
					return
				}
				if status.CurrentTailnet == nil {
					logErr("not logged in (CurrentTailnet is nil)\n")
					http.Error(w, "not logged in (CurrentTailnet is nil)", http.StatusForbidden)
					return
				}
				fqdn := opts.MachineName + "." + status.CurrentTailnet.MagicDNSSuffix
				http.Redirect(w, r, "https://"+fqdn+r.RequestURI, http.StatusMovedPermanently)
			})); err != nil {
				logErr("error serving HTTP redirect: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := http.Serve(httpListener, proxy); err != nil {
				logErr("error serving HTTP: %v\n", err)
				os.Exit(1)
			}
		}
	}()

	if opts.HTTPSMode != config.HTTPSOff {
		var httpsListener net.Listener
		tcpListener, err := s.Listen("tcp", ":443")
		if err != nil {
			logErr("error listening on port 443: %v\n", err)
			os.Exit(1)
		}
		httpsListener = tls.NewListener(tcpListener, &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				start := time.Now()
				defer func() {
					log("GetCertificate took %v\n", time.Since(start))
				}()
				return lc.GetCertificate(hello)
			},
		})
		go func() {
			defer httpsListener.Close()
			if err := http.Serve(httpsListener, proxy); err != nil {
				logErr("error serving HTTPS: %v\n", err)
				os.Exit(1)
			}
		}()
	}

	log("listening as %s, forwarding to %v\n", opts.MachineName, opts.Target)

	select {}
}
