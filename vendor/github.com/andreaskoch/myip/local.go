// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package myip

import (
	"net"
)

// NewLocalIPProvider creates a new instance of the
// LocalIPProvider type.
func NewLocalIPProvider() (LocalIPProvider, error) {
	localNetworkAddressProvider, err := newInterfaceIPProvider()
	if err != nil {
		return LocalIPProvider{}, err
	}

	return LocalIPProvider{localNetworkAddressProvider}, nil
}

// LocalIPProvider provides access to local
// IP addresses.
type LocalIPProvider struct {
	localNetworkAddressProvider IPProvider
}

// GetIPv6Addresses returns all available local IPv6 addresses.
func (p LocalIPProvider) GetIPv6Addresses() ([]net.IP, error) {

	// get the available IPs from the address provider
	allIPs, err := p.localNetworkAddressProvider.GetIPs()
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

// GetIPv4Addresses returns all local IPv4 addresses.
func (p LocalIPProvider) GetIPv4Addresses() ([]net.IP, error) {

	// get the available IPs from the address provider
	allIPs, err := p.localNetworkAddressProvider.GetIPs()
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

// newInterfaceIPProvider creates a new instance of the interfaceAddressProvider type
// with the local network interfaces as a data source.
func newInterfaceIPProvider() (interfaceAddressProvider, error) {

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
