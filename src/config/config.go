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
	FunnelMode  = FunnelOff
	MachineName string
	Target      *url.URL
	PProf       = false
	DataDir     = ""
)

var (
	help = flag.Bool("help", false, "show help")
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

	optionsMissing := loadConfigFromEnv()
	if len(optionsMissing) > 0 {
		logger.Err("info: missing environment variable: %v. Using command line flags instead.\n", optionsMissing)
		loadConfigFromCLI()
	}

	if FunnelMode != FunnelOff && HTTPSMode == HTTPSOff {
		logger.Fatal("funnel requires HTTPS, but HTTPS is disabled.\n")
		os.Exit(1)
	}
	if FunnelMode == FunnelOnly && HTTPSMode != HTTPSOnly {
		logger.Fatal("funnel only mode requires HTTPS only mode.\n")
		os.Exit(1)
	}
	if FunnelMode == FunnelOnly && PProf {
		logger.Fatal("funnel only mode does not support pprof.\n")
		os.Exit(1)
	}
}
