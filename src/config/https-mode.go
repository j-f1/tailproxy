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

func (m HTTPSModeValue) String() string {
	switch m {
	case HTTPSOff:
		return "off"
	case HTTPSRedirect:
		return "redirect"
	case HTTPSOnly:
		return "only"
	case HTTPSBoth:
		return "both"
	default:
		return fmt.Sprintf("unknown https mode %d", m)
	}
}

func parseHTTPSMode(s string) HTTPSModeValue {
	switch s {
	case "off":
		return HTTPSOff
	case "redirect":
		return HTTPSRedirect
	case "only":
		return HTTPSOnly
	case "both":
		return HTTPSBoth
	default:
		logger.Fatal("invalid https mode %q", s)
		return -1
	}
}
