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
	// --https=off (default, only serve HTTP)
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

func parseHTTPSMode(s string) (httpsMode, error) {
	switch s {
	case "off":
		return httpsOff, nil
	case "redirect":
		return httpsRedirect, nil
	case "only":
		return httpsOnly, nil
	case "both":
		return httpsBoth, nil
	default:
		return 0, fmt.Errorf("invalid https mode %q", s)
	}
}

type options struct {
	httpsMode   httpsMode
	tailnetHost string
	target      *url.URL
}

const (
	envHTTPSMode   = "TAILPROXY_HTTPS_MODE"
	envTailnetHost = "TAILPROXY_TAILNET_HOST"
	envTarget      = "TAILPROXY_TARGET"
)

func parseOptions() options {
	var opts options
	opts.httpsMode = httpsRedirect

	// env vars
	var optionsMissing []string
	var err error
	if os.Getenv(envHTTPSMode) != "" {
		opts.httpsMode, err = parseHTTPSMode(os.Getenv(envHTTPSMode))
		if err != nil {
			fmt.Fprintf(os.Stderr, "tailproxy: %v\n", err)
			os.Exit(1)
		}
	}

	if os.Getenv(envTailnetHost) != "" {
		opts.tailnetHost = os.Getenv(envTailnetHost)
	} else {
		optionsMissing = append(optionsMissing, envTailnetHost)
	}

	if os.Getenv(envTarget) != "" {
		opts.target, err = url.Parse("http://" + os.Getenv(envTarget))
		if err != nil {
			fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
			os.Exit(1)
		}
	} else {
		optionsMissing = append(optionsMissing, envTarget)
	}

	if len(optionsMissing) == 1 {
		fmt.Fprintf(os.Stderr, "tailproxy: info: missing environment variable: %v. Using command line flags instead.\n", optionsMissing)
	} else {
		return opts
	}

	// CLI flags

	flag.Usage = func() {
		fmt.Printf("usage: %s [flags] <tailnet host> <target host:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
	}

	opts.tailnetHost = flag.Arg(0)
	opts.target, err = url.Parse("http://" + flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
		flag.Usage()
	}

	opts.httpsMode, err = parseHTTPSMode(*https)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: %v\n", err)
	}

	return opts
}

func main() {
	opts := parseOptions()

	s := new(tsnet.Server)
	s.Hostname = opts.tailnetHost
	s.Ephemeral = true

	if err := s.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error starting server: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	lc, err := s.LocalClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error getting local client: %v\n", err)
		os.Exit(1)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			fmt.Printf("tailproxy: %v %v %v\n", r.In.RemoteAddr, r.In.Method, r.In.URL)
			r.SetXForwarded()
			r.SetURL(opts.target)
			r.Out.Host = r.In.Host
		},
	}

	httpListener, err := s.Listen("tcp", ":80")
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: error listening on port 80: %v\n", err)
		os.Exit(1)
	}
	go func() {
		defer httpListener.Close()
		if opts.httpsMode == httpsRedirect {
			if err := http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			})); err != nil {
				fmt.Fprintf(os.Stderr, "tailproxy: error serving HTTP redirect: %v\n", err)
				os.Exit(1)
			}
		} else {
			if err := http.Serve(httpListener, proxy); err != nil {
				fmt.Fprintf(os.Stderr, "tailproxy: error serving HTTP: %v\n", err)
				os.Exit(1)
			}
		}
	}()

	var httpsListener net.Listener
	if opts.httpsMode != httpsOff {
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

	fmt.Printf("tailproxy: listening as %s, forwarding to %v\n", opts.tailnetHost, opts.target)

	select {}
}
