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

	if FunnelMode == FunnelOnly && HTTPSMode != HTTPSOff {
		logger.Log("note: https mode is ignored in funnel only mode.\n")
	}
}
