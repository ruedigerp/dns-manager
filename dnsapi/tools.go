package dnsapi

import (
	"fmt"
	"net"
)

func IsValidIP(ip string) bool {

	return net.ParseIP(ip) != nil

}

func CheckEmpty(value, name, flag string) bool {

	if value == "" {
		fmt.Printf("Please provide %s. (%s)\n", name, flag)
		return true
	}

	return false

}
