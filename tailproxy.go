package main

import (
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/serve"
	"tailproxy/src/ts"
)

func main() {
	config.Parse()
	logger.Log("connecting as %s, forwarding to %v. HTTPS mode: %s. Storing data in: '%s'\n", config.MachineName, config.Target, config.HTTPSMode, config.DataDir)

	ts.StartServer()
	defer ts.ShutdownServer()

	if config.HTTPSMode != config.HTTPSOnly {
		go serve.ServeHTTP()
	}
	if config.HTTPSMode != config.HTTPSOff {
		go serve.ServeHTTPS()
	}
	if config.PProf {
		go serve.ServePProf()
	}
	select {}
}
