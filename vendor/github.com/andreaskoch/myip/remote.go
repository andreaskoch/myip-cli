// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package myip

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// NewRemoteIPProvider creates a new instance of the
// RemoteIPProvider type.
func NewRemoteIPProvider() (RemoteIPProvider, error) {

	// remote IPv4 provider
	ipv4AddressProvider, ipv4ProviderErr := newRemoteAddressProvider("https://ipv4.icanhazip.com")
	if ipv4ProviderErr != nil {
		return RemoteIPProvider{}, ipv4ProviderErr
	}

	// remote IPv6 provider
	ipv6AddressProvider, ipv6ProviderErr := newRemoteAddressProvider("https://ipv6.icanhazip.com")
	if ipv6ProviderErr != nil {
		return RemoteIPProvider{}, ipv6ProviderErr
	}

	return RemoteIPProvider{
		ipv4Provider: ipv4AddressProvider,
		ipv6Provider: ipv6AddressProvider,
	}, nil
}

// RemoteIPProvider provides access to remote
// IP addresses.
type RemoteIPProvider struct {
	ipv4Provider remoteAddressProvider
	ipv6Provider remoteAddressProvider
}

// GetIPv6Addresses returns the remote IPv6 address.
func (p RemoteIPProvider) GetIPv6Addresses() ([]net.IP, error) {
	ip, err := p.ipv6Provider.GetRemoteIPAddress()
	if err != nil {
		return []net.IP{}, err
	}

	if !isIPv6(ip) {
		return []net.IP{}, fmt.Errorf("The returned IP address (%s) is not an IPv6 address", ip)
	}

	return []net.IP{ip}, nil
}

// GetIPv4Addresses returns the remote IPv4 address.
func (p RemoteIPProvider) GetIPv4Addresses() ([]net.IP, error) {
	ip, err := p.ipv4Provider.GetRemoteIPAddress()
	if err != nil {
		return []net.IP{}, err
	}

	if !isIPv4(ip) {
		return []net.IP{}, fmt.Errorf("The returned IP address (%s) is not an IPv4 address", ip)
	}

	return []net.IP{ip}, nil
}

// newRemoteAddressProvider creates a new instance of the remoteAddressProvider type
// with the local interfaces as a data source.
func newRemoteAddressProvider(providerURL string) (remoteAddressProvider, error) {
	return remoteAddressProvider{
		providerURL: providerURL,
		timeout:     time.Second * 3,
	}, nil
}

// remoteAddressProvider provides functions for accessing the IP addresses of network interfaces.
type remoteAddressProvider struct {
	providerURL string
	timeout     time.Duration
}

// GetRemoteIPAddress returns the IP address returned by the provider with the given URL.
func (r remoteAddressProvider) GetRemoteIPAddress() (net.IP, error) {

	// create a http client (allow insecure SSL certs)
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: transportConfig,
		Timeout:   r.timeout,
	}

	// ask the remote service for the IP
	resp, err := httpClient.Get(r.providerURL)
	if err != nil {
		return nil, err
	}

	// read the response
	response := make([]byte, 48)
	responseReader := bufio.NewReader(resp.Body)
	bytesRead, readErr := responseReader.Read(response)
	if readErr != nil || bytesRead == 0 {
		return nil, readErr
	}

	// prepare the response for parsing
	content := strings.TrimSpace(string(response[:bytesRead]))

	// parse the response
	ip := net.ParseIP(content)
	if ip == nil {
		return nil, fmt.Errorf("%q is not a valid IP address", content)
	}

	return ip, nil
}
