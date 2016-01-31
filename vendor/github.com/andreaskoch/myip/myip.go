// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package myip implements access to local and remote
// IP addresses. It aims to provide a common interface
// for determining the available local and remote
// IPv6 and IPv4 addresses.
package myip

import (
	"net"
)

// AddressProvider returns IP addresses.
type AddressProvider interface {
	GetIPs() (ips []net.IP, err error)
}

// The IPAddresser interface provides functions for
// retrieving IPv4 and IPv6 addresses.
type IPAddresser interface {
	IPv4Addresser
	IPv6Addresser
}

// The IPv6Addresser interface provides functions for
// retrieving IPv6 addresses.
type IPv6Addresser interface {
	GetIPv6Addresses() ([]net.IP, error)
}

// The IPv4Addresser interface provides functions for
// retrieving IPv4 addresses.
type IPv4Addresser interface {
	GetIPv4Addresses() ([]net.IP, error)
}

// The IPProvider interface returns IP addresses from a data source.
type IPProvider interface {
	// GetIPs returns all IPs available to this provider or an error
	// if the IPs cannot be accessed.
	GetIPs() (ips []net.IP, err error)
}
