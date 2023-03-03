package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"tailproxy/src/logger"
)

var (
	HTTPSMode   = HTTPSOff
	MachineName string
	Target      *url.URL
	PProf       = false
)

const (
	HTTPSOff HTTPSModeValue = iota
	HTTPSRedirect
	HTTPSOnly
	HTTPSBoth
)

var (
	// --https=off (default, only serve HTTP)
	// --https=redirect (redirect HTTP to HTTPS)
	// --https=only (only serve HTTPS)
	// --https=both (serve both HTTP and HTTPS)
	https = flag.String("https", "off", "HTTPS mode (off, on, only, both)")

	pprof = flag.Bool("pprof", false, "enable pprof")

	help = flag.Bool("help", false, "show help")
)

type HTTPSModeValue int

func (m HTTPSModeValue) String() string {
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

func parseHTTPSMode(s string) HTTPSModeValue {
	switch s {
	case "off":
		return HTTPSOff
	case "redirect":
		return HTTPSRedirect
	case "only":
		return HTTPSOnly
	case "both":
		return HTTPSBoth
	default:
		logger.Fatal("invalid https mode %q", s)
		return -1
	}
}

const (
	envHTTPSMode = "TAILPROXY_HTTPS_MODE"
	envName      = "TAILPROXY_NAME"
	envTarget    = "TAILPROXY_TARGET"
	envPProf     = "TAILPROXY_PPROF_ENABLED"
)

func Parse() {
	flag.Usage = func() {
		fmt.Printf("usage: %s [flags] <tailnet host> <target host:port>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if *help {
		flag.Usage()
	}

	// env vars
	var optionsMissing []string
	var err error
	if os.Getenv(envHTTPSMode) != "" {
		HTTPSMode = parseHTTPSMode(os.Getenv(envHTTPSMode))
	}

	if os.Getenv((envPProf)) != "" {
		PProf = true
	}

	if os.Getenv(envName) != "" {
		MachineName = os.Getenv(envName)
	} else {
		optionsMissing = append(optionsMissing, envName)
	}

	if os.Getenv(envTarget) != "" {
		Target, err = url.Parse("http://" + os.Getenv(envTarget))
		if err != nil {
			logger.Fatal("invalid target: %v\n", err)
		}
	} else {
		optionsMissing = append(optionsMissing, envTarget)
	}

	if len(optionsMissing) > 0 {
		logger.Err("info: missing environment variable: %v. Using command line flags instead.\n", optionsMissing)
		// CLI flags

		if flag.NArg() != 2 {
			flag.Usage()
		}

		MachineName = flag.Arg(0)
		Target, err = url.Parse("http://" + flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
			flag.Usage()
		}

		PProf = *pprof
		HTTPSMode = parseHTTPSMode(*https)
	}
}
