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

const timeout = 10

// NewRemoteIPProvider creates a new instance of the
// RemoteIPProvider type.
func NewRemoteIPProvider() RemoteIPProvider {

	return RemoteIPProvider{
		ipv4Providers: []remoteAddressProvider{
			newRemoteIPv4AddressProvider("https://ipv4.yip.li"),
			newRemoteIPv4AddressProvider("https://ipv4.icanhazip.com"),
		},
		ipv6Providers: []remoteAddressProvider{
			newRemoteIPv6AddressProvider("https://ipv6.yip.li"),
			newRemoteIPv6AddressProvider("https://ipv6.icanhazip.com"),
		},
	}

}

// RemoteIPProvider provides access to remote
// IP addresses.
type RemoteIPProvider struct {
	ipv4Providers []remoteAddressProvider
	ipv6Providers []remoteAddressProvider
}

// GetIPv6Addresses returns the remote IPv6 address.
func (p RemoteIPProvider) GetIPv6Addresses() ([]net.IP, error) {

	ip, err := requestRemoteIP(p.ipv6Providers)
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

	ip, err := requestRemoteIP(p.ipv4Providers)
	if err != nil {
		return []net.IP{}, err
	}

	if !isIPv4(ip) {
		return []net.IP{}, fmt.Errorf("The returned IP address (%s) is not an IPv4 address", ip)
	}

	return []net.IP{ip}, nil
}

func requestRemoteIP(providers []remoteAddressProvider) (net.IP, error) {

	if len(providers) == 0 {
		return nil, fmt.Errorf("No providers given")
	}

	numberOfProviders := len(providers)

	ips := make(chan net.IP, numberOfProviders)

	for _, provider := range providers {

		currentProvider := provider

		go func() {

			ip, _ := currentProvider.GetRemoteIPAddress()
			ips <- ip

		}()

	}

	for {
		select {
		case ip := <-ips:
			{
				if ip != nil {
					return ip, nil
				}
			}
		case <-time.After(time.Second * timeout):
			return nil, fmt.Errorf("Timeout")
		}
	}
}

// newRemoteIPv4AddressProvider creates a new instance of the remoteAddressProvider type
// with the given provider URL as the data source over IPv4.
func newRemoteIPv4AddressProvider(providerURL string) remoteAddressProvider {
	return newRemoteAddressProvider("tcp4", providerURL)
}

// newRemoteIPv6AddressProvider creates a new instance of the remoteAddressProvider type
// with the given provider URL as the data source over IPv4.
func newRemoteIPv6AddressProvider(providerURL string) remoteAddressProvider {
	return newRemoteAddressProvider("tcp6", providerURL)
}

// newRemoteAddressProvider creates a new instance of the remoteAddressProvider type
// with the given provider URL as the data source over the given network ("tcp", "tcp6", "tcp4")
func newRemoteAddressProvider(network, providerURL string) remoteAddressProvider {
	return remoteAddressProvider{
		network:     network,
		providerURL: providerURL,
		timeout:     time.Second * timeout,
	}
}

// remoteAddressProvider provides functions for accessing the IP addresses of network interfaces.
type remoteAddressProvider struct {
	network     string
	providerURL string
	timeout     time.Duration
}

// GetRemoteIPAddress returns the IP address returned by the provider with the given URL.
func (r remoteAddressProvider) GetRemoteIPAddress() (net.IP, error) {

	// create a http client (allow insecure SSL certs)
	dialer := func(network, address string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout: r.timeout,
		}
		return dialer.Dial(r.network, address)
	}

	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            dialer,
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
