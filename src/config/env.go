package config

import (
	"net/url"
	"os"
	"tailproxy/src/logger"
)

const (
	envHTTPSMode = "TAILPROXY_HTTPS_MODE"
	envName      = "TAILPROXY_NAME"
	envTarget    = "TAILPROXY_TARGET"
	envPProf     = "TAILPROXY_PPROF_ENABLED"
)

func loadConfigFromEnv() []string {
	var optionsMissing []string
	if os.Getenv(envHTTPSMode) != "" {
		HTTPSMode = parseHTTPSMode(os.Getenv(envHTTPSMode))
	}

	if os.Getenv(envPProf) != "" {
		PProf = true
	}

	if os.Getenv(envName) != "" {
		MachineName = os.Getenv(envName)
	} else {
		optionsMissing = append(optionsMissing, envName)
	}

	if os.Getenv(envTarget) != "" {
		var err error
		Target, err = url.Parse("http://" + os.Getenv(envTarget))
		if err != nil {
			logger.Fatal("invalid target: %v\n", err)
		}
	} else {
		optionsMissing = append(optionsMissing, envTarget)
	}

	return optionsMissing
}
