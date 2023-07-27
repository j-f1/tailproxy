package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
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
		logger.Err("info: missing environment variable: %v. Using command line flags instead.", optionsMissing)
		loadConfigFromCLI()
	}

	if (FunnelMode == FunnelOnly || FunnelMode == FunnelRedirect) && HTTPSMode != HTTPSOff {
		logger.Log("note: HTTPS mode is ignored in Funnel-only mode.")
	}

	if DataDir == "" {
		if _, err := os.Stat("/data"); err == nil {
			DataDir = "/data"
		}
	}
}

func parseTarget(target string) {
	if !strings.Contains(target, "://") {
		target = "http://" + target
	}
	var err error
	Target, err = url.Parse(target)
	if err != nil {
		logger.Fatal("invalid target: %v", err)
	}
}
