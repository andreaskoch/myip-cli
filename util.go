// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
)

// isIPv4 returns true if the given ip address is an IPv4 address.
func isIPv4(ip net.IP) bool {
	if ip.To4() != nil {
		return true
	}
	return false
}

// isIPv6 returns true if the given ip address is an IPv6 address.
func isIPv6(ip net.IP) bool {
	return isIPv4(ip) == false
}

// getIP returns the IP of the given address.
func getIP(address net.Addr) net.IP {

	var ip net.IP
	switch v := address.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}

	return ip
}

func isLoopbackIP(ip net.IP) bool {

	// ignore loopback ips
	if ip.IsLoopback() {
		return true
	}

	// ignore local ips
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	return false
}
