// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// useIPv4 contains a flag inidicating whether IPv4 addresses should be used (default: false)
var useIPv4 bool

func init() {

	executableName := os.Args[0]

	flag.BoolVar(&useIPv4, "4", false, fmt.Sprintf("Forces %s to use IPv4 addresses only.", executableName))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s returns your local IPv6 (or IPv4) address.\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
}

// mylocalip detects the current local IPv6 (or IPv4) address and prints it to stdout
func main() {
	flag.Parse()

	ips, err := getIPs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

	for _, ip := range ips {

		// IPv6
		if isIPv6(ip) && useIPv4 == false {
			fmt.Fprintf(os.Stdout, "%s", ip)
			os.Exit(0)
		}

		// IPv4
		if isIPv4(ip) && useIPv4 == true {
			fmt.Fprintf(os.Stdout, "%s", ip)
			os.Exit(0)
		}
	}

	fmt.Fprintf(os.Stderr, "No addresses detected")
	os.Exit(1)
}

// getIPs returns all IP addresses of the current machine.
func getIPs() (ips []net.IP, err error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			ip := getIP(addr)

			// ignore loopback ips
			if ip.IsLoopback() {
				continue
			}

			// ignore local ips
			if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}

			ips = append(ips, ip)
		}
	}

	return ips, nil
}

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
