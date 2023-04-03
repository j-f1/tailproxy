package config

import (
	"fmt"
	"tailproxy/src/logger"
)

type HTTPSModeValue int

const (
	// Default, only serve HTTP
	HTTPSOff HTTPSModeValue = iota
	// Redirect HTTP to HTTPS
	HTTPSRedirect
	// Only serve HTTPS
	HTTPSOnly
	// Serve both HTTP and HTTPS
	HTTPSBoth
)

const (
	httpsOff      = "off"
	httpsRedirect = "redirect"
	httpsOnly     = "only"
	httpsBoth     = "both"
)

func (m HTTPSModeValue) String() string {
	switch m {
	case HTTPSOff:
		return httpsOff
	case HTTPSRedirect:
		return httpsRedirect
	case HTTPSOnly:
		return httpsOnly
	case HTTPSBoth:
		return httpsBoth
	default:
		return fmt.Sprintf("unknown https mode %d", m)
	}
}

func parseHTTPSMode(s string) HTTPSModeValue {
	switch s {
	case httpsOff:
		return HTTPSOff
	case httpsRedirect:
		return HTTPSRedirect
	case httpsOnly:
		return HTTPSOnly
	case httpsBoth:
		return HTTPSBoth
	default:
		logger.Fatal("unknown https mode %q", s)
		return -1
	}
}
