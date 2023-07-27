package config

import (
	"fmt"
	"tailproxy/src/logger"
)

type FunnelModeValue int

const (
	// Default, only serve on the tailnet
	FunnelOff FunnelModeValue = iota
	// serve to both tailnet and funnel
	FunnelOn
	// only serve to funnel
	FunnelOnly
	// redirect traffic to funnel
	FunnelRedirect
)

const (
	funnelOff      = "off"
	funnelOn       = "on"
	funnelOnly     = "only"
	funnelRedirect = "redirect"
)

func (m FunnelModeValue) String() string {
	switch m {
	case FunnelOff:
		return funnelOff
	case FunnelOn:
		return funnelOn
	case FunnelOnly:
		return funnelOnly
	case FunnelRedirect:
		return funnelRedirect
	default:
		return fmt.Sprintf("unknown funnel mode %d", m)
	}
}

func parseFunnelMode(s string) FunnelModeValue {
	switch s {
	case funnelOff:
		return FunnelOff
	case funnelOn:
		return FunnelOn
	case funnelOnly:
		return FunnelOnly
	case funnelRedirect:
		return FunnelRedirect
	default:
		logger.Fatal("invalid funnel mode %q", s)
		return -1
	}
}
