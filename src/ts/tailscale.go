package ts

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"tailproxy/src/config"
	"tailproxy/src/logger"

	"tailscale.com/client/tailscale"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tsnet"
)

var s = new(tsnet.Server)
var lc *tailscale.LocalClient

func StartServer() {
	s.Hostname = config.MachineName
	s.Ephemeral = true

	var err error
	lc, err = s.LocalClient()
	if err != nil {
		logger.Fatal("error starting server: %v\n", err)
	}

	err = lc.StartLoginInteractive(context.Background())
	if err != nil {
		logger.Fatal("error starting login: %v\n", err)
	}
}

func Status() *ipnstate.Status {
	status, err := lc.Status(context.Background())
	if err != nil || status == nil {
		logger.Fatal("error getting profile status: %v\n", err)
	}
	return status
}

func ShutdownServer() {
	s.Close()
}

func Listen(network, address string) net.Listener {
	listener, err := s.Listen(network, address)
	if err != nil {
		logger.Fatal("error listening for %s on %s: %v\n", network, address, err)
	}
	return listener
}

func MagicDNSSuffix() (string, string) {
	status, err := lc.Status(context.Background())
	if err != nil || status == nil {
		logger.Err("error getting profile status: %v\n", err)
		return "", fmt.Sprintf("error getting profile status: %v", err)
	}
	if status.CurrentTailnet == nil {
		logger.Err("not logged in (CurrentTailnet is nil)\n")
		return "", "not logged in (CurrentTailnet is nil)"
	}
	return status.CurrentTailnet.MagicDNSSuffix, ""
}

func GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return lc.GetCertificate(hello)
}
