package dnsapi

import (
	"net"
)

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
