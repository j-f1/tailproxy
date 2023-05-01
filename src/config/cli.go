package config

import (
	"flag"
)

var (
	https   = flag.String("https", "off", "HTTPS mode (off, on, only, both)")
	funnel  = flag.String("funnel", "off", "funnel mode (off, on, only)")
	dataDir = flag.String("data-dir", "", "data directory")
	pprof   = flag.Bool("pprof", false, "enable pprof")
)

func loadConfigFromCLI() {
	if flag.NArg() != 2 {
		flag.Usage()
	}

	MachineName = flag.Arg(0)
	parseTarget(flag.Arg(1))

	HTTPSMode = parseHTTPSMode(*https)
	FunnelMode = parseFunnelMode(*funnel)
	DataDir = *dataDir
	PProf = *pprof
}
