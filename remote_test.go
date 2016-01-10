// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_GetIPv4Address_ProviderReturnsValidIPv4Address_IPReturned_NoErrorIsReturned(t *testing.T) {
	// arrange
	response := "192.168.22.1"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	remoteIPAddressProviderURL := testServer.URL

	remoteAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)

	// act
	remoteIPProvider := remoteIPProvider{ipv4Provider: remoteAddressProvider}
	ips, err := remoteIPProvider.GetIPv4Addresses()

	// assert
	if len(ips) == 0 {
		t.Fail()
		t.Logf("GetIPv4Address() should return IPs if the IP provider responded with %q.", response)
	}

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetIPv4Address() returned an error even though there should be no reason for it. %q is a valid IPv4 address.", response)
	}
}

func Test_GetIPv4Address_ProviderReturnsIPv6Address_NoIPsReturned_ErrorIsReturned(t *testing.T) {
	// arrange
	response := "2001:0db8:0000:0042:0000:8a2e:0370:7334"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	remoteIPAddressProviderURL := testServer.URL

	remoteAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)

	// act
	remoteIPProvider := remoteIPProvider{ipv4Provider: remoteAddressProvider}
	ips, err := remoteIPProvider.GetIPv4Addresses()

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("GetIPv4Address() should return nil if the IP provider responded with something that is not an IPv4 address (%q).", response)
	}

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetIPv4Address() returned no error even though the provider returned an invalid IPv4 address: %q", response)
	}
}

func Test_GetIPv6Address_ProviderReturnsValidIPv6Address_IPsAreReturned_NoErrorIsReturned(t *testing.T) {
	// arrange
	response := "2001:0db8:0000:0042:0000:8a2e:0370:7334"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	remoteIPAddressProviderURL := testServer.URL

	remoteAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)

	// act
	remoteIPProvider := remoteIPProvider{ipv6Provider: remoteAddressProvider}
	ips, err := remoteIPProvider.GetIPv6Addresses()

	// assert
	if len(ips) == 0 {
		t.Fail()
		t.Logf("GetIPv6Address() should not return nil if the IP provider responded with %q.", response)
	}

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetIPv6Address() returned an error even though there should be no reason for it. %q is a valid IPv4 address.", response)
	}
}

func Test_GetIPv6Address_ProviderReturnsIPv4Address_NoIPsReturned_ErrorIsReturned(t *testing.T) {
	// arrange
	response := "192.168.22.1"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	remoteIPAddressProviderURL := testServer.URL

	remoteAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)

	// act
	remoteIPProvider := remoteIPProvider{ipv6Provider: remoteAddressProvider}
	ips, err := remoteIPProvider.GetIPv6Addresses()

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("GetIPv6Address() should not return any IPs if the IP provider responded with something that is not an IPv6 address (%q).", response)
	}

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetIPv6Address() returned no error even though the provider returned an invalid IPv6 address: %q", response)
	}
}

// If a valid URL is provided and the response is a valid IP getRemoteIPAddress should not return an error.
func Test_getRemoteIPAddress_ValidProviderURL_IPIsReturned(t *testing.T) {
	// arrange
	response := "192.168.22.1"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer testServer.Close()
	remoteIPAddressProviderURL := testServer.URL

	// act
	ipAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)
	ip, err := ipAddressProvider.GetRemoteIPAddress()

	// assert
	if err != nil {
		t.Fail()
		t.Logf("getRemoteIPAddress() returned an error: %s", err.Error())
	}

	if ip == nil {
		t.Fail()
		t.Logf("getRemoteIPAddress() did not return an IP address.")
	}
}

func Test_getRemoteIPAddress_InvalidProviderURL_RequestTimesOut(t *testing.T) {
	// arrange
	expectedTimeout := 3 * time.Second
	response := "192.168.22.1"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		fmt.Fprintln(w, response)
	}))
	defer testServer.Close() // stop the server
	remoteIPAddressProviderURL := testServer.URL

	// act
	startTime := time.Now()
	ipAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)
	ipAddressProvider.GetRemoteIPAddress()
	elapsed := time.Since(startTime)

	// assert
	elapsedSeconds := int(elapsed.Seconds())
	expectedSeconds := int(expectedTimeout.Seconds())
	if elapsedSeconds > expectedSeconds {
		t.Fail()
		t.Logf("getRemoteIPAddress() did not time out as expected. Expected timout: %v, Time elapsed: %v", expectedSeconds, elapsedSeconds)
	}
}

// If the supplied remote IP provider URL is invalid or does not exist getRemoteIPAddress should return an error.
func Test_getRemoteIPAddress_InvalidProviderURL_ErrorIsReturned(t *testing.T) {
	// arrange
	response := "192.168.22.1"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	remoteIPAddressProviderURL := testServer.URL

	testServer.Close() // stop the server

	// act
	ipAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)
	_, err := ipAddressProvider.GetRemoteIPAddress()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("getRemoteIPAddress() returned no error - even though the supplied URL (%s) was invalid.", remoteIPAddressProviderURL)
	}
}

// If the remote ip server returns gibberish getRemoteIPAddress should return an error.
func Test_getRemoteIPAddress_ValidProviderURL_ResponseContentIsInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	invalidResponses := []string{
		"",
		"   ",
		"aaa",
		"192.168.1.1.1",
		"192.168.1.260",
		"2001 0db8 0000 0042 0000 8a2e 0370 7334",
		"2001 : 0db8 : 0000 : 0042 : 0000 : 8a2e : 0370 : 7334",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334:0000:1111",
	}

	for _, response := range invalidResponses {

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, response)
		}))
		defer testServer.Close()
		remoteIPAddressProviderURL := testServer.URL

		// act
		ipAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)
		_, err := ipAddressProvider.GetRemoteIPAddress()

		// assert
		if err == nil {
			t.Fail()
			t.Logf("getRemoteIPAddress() should return an error if the IP provider returns an invalid response (%q).", response)
		}

	}
}

// If the remote ip server returns gibberish getRemoteIPAddress should return an error.
func Test_getRemoteIPAddress_ValidProviderURL_ResponseContentIsValid_IPIsNotNil(t *testing.T) {
	// arrange
	invalidResponses := []string{
		"0.0.0.0",
		"1.1.1.1",
		"255.255.255.255",
		"::1",
		"::",
		"0:0:0:0:0:0:0:0",
		"2001::0370:7334",
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		"  192.169.12.2   ",
		"  2001:0db8:0000:0042:0000:8a2e:0370:7334 ",
	}

	for _, response := range invalidResponses {

		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, response)
		}))
		defer testServer.Close()
		remoteIPAddressProviderURL := testServer.URL

		// act
		ipAddressProvider, _ := newRemoteAddressProvider(remoteIPAddressProviderURL)
		ip, _ := ipAddressProvider.GetRemoteIPAddress()

		// assert
		if ip == nil {
			t.Fail()
			t.Logf("getRemoteIPAddress() should return an IP if the response is %q but returned nil instead.", response)
		}

	}
}
