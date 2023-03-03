package main

import (
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/serve"
	"tailproxy/src/ts"
)

func main() {
	opts := config.ParseOptions()

	logger.Log("machine name: %v, target: %v, https mode: %v\n", opts.MachineName, opts.Target, opts.HTTPSMode)

	ts.StartServer(opts)
	defer ts.ShutdownServer()

	proxy := serve.MakeProxy(opts)

	go serve.ServeHTTP(opts, proxy)

	if opts.HTTPSMode != config.HTTPSOff {
		go serve.ServeHTTPS(opts, proxy)
	}

	logger.Log("listening as %s, forwarding to %v\n", opts.MachineName, opts.Target)

	select {}
}
