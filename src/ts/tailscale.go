package ts

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"tailproxy/src/config"
	"tailproxy/src/logger"

	"tailscale.com/client/tailscale"
	"tailscale.com/client/tailscale/apitype"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tsnet"
)

var s *tsnet.Server
var lc *tailscale.LocalClient

func StartServer() {
	s = &tsnet.Server{
		Hostname:  config.MachineName,
		Ephemeral: true,
	}
	if len(config.DataDir) > 0 {
		err := os.MkdirAll(path.Join(config.DataDir, "tailscale"), 0700)
		if err != nil {
			logger.Fatal("error creating data dir: %v", err)
		}
		s.Dir = path.Join(config.DataDir, "tailscale")
	}

	var err error
	lc, err = s.LocalClient()
	if err != nil {
		logger.Fatal("error getting local client: %v", err)
	}

	err = lc.StartLoginInteractive(context.Background())
	if err != nil {
		logger.Fatal("error starting login: %v", err)
	}
}

func Status() *ipnstate.Status {
	status, err := lc.Status(context.Background())
	if err != nil || status == nil {
		logger.Fatal("error getting status: %v", err)
	}
	return status
}

func ShutdownServer() {
	s.Close()
	s = nil
}

func ListenTailnet(port int) net.Listener {
	return listen(port, false)
}

func ListenFunnel(port int) net.Listener {
	return listen(port, true)
}

func listen(port int, onFunnel bool) net.Listener {
	addr := fmt.Sprintf(":%d", port)
	network := "tcp"

	if onFunnel {
		listener, err := s.ListenFunnel(network, addr, tsnet.FunnelOnly())
		if err != nil {
			logger.Fatal("error listening for %s on port %v (funnel): %v", network, port, err)
		}
		return listener
	} else {
		listener, err := s.Listen(network, addr)
		if err != nil {
			logger.Fatal("error listening for %s on port %v: %v", network, port, err)
		}
		return listener
	}

}

func MagicDNSSuffix(ctx context.Context) (string, string) {
	status, err := lc.Status(ctx)
	if err != nil || status == nil {
		logger.Err("error getting profile status: %v", err)
		return "", fmt.Sprintf("error getting profile status: %v", err)
	}
	if status.CurrentTailnet == nil {
		logger.Err("not logged in (CurrentTailnet is nil)")
		return "", "not logged in (CurrentTailnet is nil)"
	}
	return status.CurrentTailnet.MagicDNSSuffix, ""
}

func WhoIs(r *http.Request) (*apitype.WhoIsResponse, error) {
	return lc.WhoIs(r.Context(), r.RemoteAddr)
}

func GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return lc.GetCertificate(hello)
}
