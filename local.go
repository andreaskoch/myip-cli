// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
)

// newLocalIPProvider creates a new instance of the
// localIPProvider type.
func newLocalIPProvider() (localIPProvider, error) {
	addressProvider, err := newInterfaceAddressProvider()
	if err != nil {
		return localIPProvider{}, err
	}

	return localIPProvider{addressProvider}, nil
}

// localIPProvider provides access to local
// IP addresses.
type localIPProvider struct {
	addressProvider interfaceAddressProvider
}

// GetIPv6Address returns the first local IPv6 address
func (p localIPProvider) GetIPv6Addresses() ([]net.IP, error) {

	// get the available IPs from the address provider
	allIPs, err := p.addressProvider.GetIPs()
	if err != nil {
		return []net.IP{}, err
	}

	var filteredIPs []net.IP
	for _, ip := range allIPs {

		// ignore loopback IPs
		if isLoopbackIP(ip) {
			continue
		}

		// ignore all non-IPv6 addresses
		if !isIPv6(ip) {
			continue
		}

		filteredIPs = append(filteredIPs, ip)
	}

	return filteredIPs, nil
}

// GetIPv4Address returns the first local IPv4 address
func (p localIPProvider) GetIPv4Addresses() ([]net.IP, error) {

	// get the available IPs from the address provider
	allIPs, err := p.addressProvider.GetIPs()
	if err != nil {
		return []net.IP{}, err
	}

	var filteredIPs []net.IP
	for _, ip := range allIPs {

		// ignore loopback IPs
		if isLoopbackIP(ip) {
			continue
		}

		// ignore all non-IPv4 addresses
		if !isIPv4(ip) {
			continue
		}

		filteredIPs = append(filteredIPs, ip)
	}

	return filteredIPs, nil
}

// newInterfaceAddressProvider creates a new instance of the interfaceAddressProvider type
// with the local interfaces as a data source.
func newInterfaceAddressProvider() (interfaceAddressProvider, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return interfaceAddressProvider{}, err
	}

	return interfaceAddressProvider{
		interfaces: interfaces,
	}, nil
}

// interfaceAddressProvider provides functions for accessing the IP addresses of network interfaces.
type interfaceAddressProvider struct {
	interfaces []net.Interface
}

// GetIPs returns all IP addresses of the current machine.
func (p interfaceAddressProvider) GetIPs() (ips []net.IP, err error) {

	for _, i := range p.interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			ip := getIP(addr)
			ips = append(ips, ip)
		}
	}

	return ips, nil
}
