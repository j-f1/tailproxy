package tailproxy

import (
	"tailproxy/src/config"
	"tailproxy/src/logger"
	"tailproxy/src/serve"
	"tailproxy/src/ts"
)

func main() {
	config.Parse()
	logger.Log("connecting as %s, forwarding to %v. HTTPS mode: %s\n", config.MachineName, config.Target, config.HTTPSMode)

	ts.StartServer()
	defer ts.ShutdownServer()

	proxy := serve.MakeProxy()

	if config.HTTPSMode != config.HTTPSOnly {
		go serve.ServeHTTP(proxy)
	}
	if config.HTTPSMode != config.HTTPSOff {
		go serve.ServeHTTPS(proxy)
	}
	select {}
}
