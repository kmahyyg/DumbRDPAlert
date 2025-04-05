package utils

import (
	"errors"
	netif "tailscale.com/net/netmon"
)

var (
	ErrCannotGetLocalIP = errors.New("cannot get local ip")
)

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() ([]string, error) {
	gLogger, err := GetLoggerInstance()
	if err != nil {
		return nil, err
	}
	dIf, _, err := netif.LocalAddresses()
	if err != nil {
		gLogger.Critical("get local addr:", err)
		return nil, err
	}
	if len(dIf) == 0 {
		return nil, ErrCannotGetLocalIP
	}
	res := make([]string, 0)
	for _, v := range dIf {
		res = append(res, v.String())
	}
	return res, nil
}
