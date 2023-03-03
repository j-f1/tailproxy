package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

var (
	// --https=off (default, only serve HTTP)
	// --https=redirect (redirect HTTP to HTTPS)
	// --https=only (only serve HTTPS)
	// --https=both (serve both HTTP and HTTPS)
	https = flag.String("https", "off", "HTTPS mode (off, on, only, both)")

	help = flag.Bool("help", false, "show help")
)

type HTTPSMode int

const (
	HTTPSOff HTTPSMode = iota
	HTTPSRedirect
	HTTPSOnly
	HTTPSBoth
)

func parseHTTPSMode(s string) (HTTPSMode, error) {
	switch s {
	case "off":
		return HTTPSOff, nil
	case "redirect":
		return HTTPSRedirect, nil
	case "only":
		return HTTPSOnly, nil
	case "both":
		return HTTPSBoth, nil
	default:
		return 0, fmt.Errorf("invalid https mode %q", s)
	}
}
func (m HTTPSMode) String() string {
	switch m {
	case HTTPSOff:
		return "off"
	case HTTPSRedirect:
		return "redirect"
	case HTTPSOnly:
		return "only"
	case HTTPSBoth:
		return "both"
	default:
		return fmt.Sprintf("unknown https mode %d", m)
	}
}

type Options struct {
	HTTPSMode   HTTPSMode
	MachineName string
	Target      *url.URL
}

const (
	envHTTPSMode = "TAILPROXY_HTTPS_MODE"
	envName      = "TAILPROXY_NAME"
	envTarget    = "TAILPROXY_TARGET"
)

func ParseOptions() Options {
	flag.Usage = func() {
		fmt.Printf("usage: %s [flags] <tailnet host> <target host:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if *help {
		flag.Usage()
	}

	var opts Options
	opts.HTTPSMode = HTTPSOff

	// env vars
	var optionsMissing []string
	var err error
	if os.Getenv(envHTTPSMode) != "" {
		opts.HTTPSMode, err = parseHTTPSMode(os.Getenv(envHTTPSMode))
		if err != nil {
			fmt.Fprintf(os.Stderr, "tailproxy: %v\n", err)
			os.Exit(1)
		}
	}

	if os.Getenv(envName) != "" {
		opts.MachineName = os.Getenv(envName)
	} else {
		optionsMissing = append(optionsMissing, envName)
	}

	if os.Getenv(envTarget) != "" {
		opts.Target, err = url.Parse("http://" + os.Getenv(envTarget))
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

	if flag.NArg() != 2 {
		flag.Usage()
	}

	opts.MachineName = flag.Arg(0)
	opts.Target, err = url.Parse("http://" + flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
		flag.Usage()
	}

	opts.HTTPSMode, err = parseHTTPSMode(*https)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: %v\n", err)
	}

	return opts
}
