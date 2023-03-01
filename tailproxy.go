package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"tailscale.com/tsnet"
)

var (
	// --https=off (default)
	// --https=redirect (redirect HTTP to HTTPS)
	// --https=only (only serve HTTPS)
	// --https=both (serve both HTTP and HTTPS)
	https = flag.String("https", "off", "HTTPS mode (off, on, only, both)")
)

type httpsMode int

const (
	httpsOff httpsMode = iota
	httpsRedirect
	httpsOnly
	httpsBoth
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [flags] <tailnet host> <target host:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
	}

	tailnetHost := flag.Arg(0)
	target, err := url.Parse("http://" + flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
		flag.Usage()
	}

	var mode httpsMode
	switch *https {
	case "off":
		mode = httpsOff
	case "on":
		mode = httpsRedirect
	case "only":
		mode = httpsOnly
	case "both":
		mode = httpsBoth
	default:
		flag.Usage()
	}

	s := new(tsnet.Server)
	s.Hostname = tailnetHost
	defer s.Close()

	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error starting server: %v\n", err)
		os.Exit(1)
	}

	lc, err := s.LocalClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error getting local client: %v\n", err)
		os.Exit(1)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetXForwarded()
			r.SetURL(target)
			r.Out.Host = r.In.Host
		},
	}

	httpListener, err := s.Listen("tcp", ":80")
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error listening on port 80: %v", err)
		os.Exit(1)
	}
	go func() {
		defer httpListener.Close()
		if mode == httpsRedirect {
			if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			})); err != nil {
				fmt.Fprintf(os.Stderr, "tailproxy: error serving HTTP redirect: %v", err)
				os.Exit(1)
			}
		} else {
			if err := http.Serve(httpListener, proxy); err != nil {
				fmt.Fprintf(os.Stderr, "tailproxy: error serving HTTP: %v", err)
				os.Exit(1)
			}
		}
	}()

	var httpsListener net.Listener
	if mode != httpsOff {
		tcpListener, err := s.Listen("tcp", ":443")
		if err != nil {
			fmt.Fprintf(os.Stderr, "tailproxy: error listening on port 443: %v\n", err)
			os.Exit(1)
		}
		httpsListener = tls.NewListener(tcpListener, &tls.Config{
			GetCertificate: lc.GetCertificate,
		})
		go func() {
			defer httpsListener.Close()
			if err := http.Serve(httpsListener, proxy); err != nil {
				fmt.Fprintf(os.Stderr, "tailproxy: error serving HTTPS: %v", err)
				os.Exit(1)
			}
		}()
	}

	fmt.Printf("tailproxy: listening as %s, forwarding to %s\n", tailnetHost, target)

	select {}
}
