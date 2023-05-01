package config

import (
	"os"
)

const (
	envHTTPSMode = "TAILPROXY_HTTPS_MODE"
	envFunnel    = "TAILPROXY_FUNNEL_MODE"
	envName      = "TAILPROXY_NAME"
	envTarget    = "TAILPROXY_TARGET"
	envPProf     = "TAILPROXY_PPROF_ENABLED"
	envDataDir   = "TAILPROXY_DATA_DIR"
)

func loadConfigFromEnv() []string {
	var optionsMissing []string
	if os.Getenv(envHTTPSMode) != "" {
		HTTPSMode = parseHTTPSMode(os.Getenv(envHTTPSMode))
	}

	if os.Getenv(envFunnel) != "" {
		FunnelMode = parseFunnelMode(os.Getenv(envFunnel))
	}

	if os.Getenv(envPProf) != "" {
		PProf = true
	}

	if os.Getenv(envDataDir) != "" {
		DataDir = os.Getenv(envDataDir)
	}

	if os.Getenv(envName) != "" {
		MachineName = os.Getenv(envName)
	} else {
		optionsMissing = append(optionsMissing, envName)
	}

	if os.Getenv(envTarget) != "" {
		parseTarget(os.Getenv(envTarget))
	} else {
		optionsMissing = append(optionsMissing, envTarget)
	}

	return optionsMissing
}
