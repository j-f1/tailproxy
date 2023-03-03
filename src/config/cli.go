package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

var (
	https   = flag.String("https", "off", "HTTPS mode (off, on, only, both)")
	dataDir = flag.String("data-dir", "", "data directory")
	pprof   = flag.Bool("pprof", false, "enable pprof")
)

func loadConfigFromCLI() {
	if flag.NArg() != 2 {
		flag.Usage()
	}

	MachineName = flag.Arg(0)
	var err error
	Target, err = url.Parse("http://" + flag.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "tailproxy: invalid target: %v\n", err)
		flag.Usage()
	}

	HTTPSMode = parseHTTPSMode(*https)
	DataDir = *dataDir
	PProf = *pprof
}
