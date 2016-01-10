// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"testing"
)

type testIPProvider struct {
	ipv4IPs []net.IP
	ipv4Err error

	ipv6IPs []net.IP
	ipv6Err error
}

func (p testIPProvider) GetIPv6Addresses() ([]net.IP, error) {
	return p.ipv6IPs, p.ipv6Err
}

func (p testIPProvider) GetIPv4Addresses() ([]net.IP, error) {
	return p.ipv4IPs, p.ipv4Err
}

// getMyIP should return the IPv4 addresses of the IP provider if the useIPv4 flag is set to true.
func Test_getMyIP_UseIPv4True_IPProviderHasIPv4Addresses_IPv4AddressesAreReturned(t *testing.T) {
	// arrange
	ipProvider := testIPProvider{
		ipv4IPs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		ipv4Err: nil,
		ipv6IPs: []net.IP{
			net.ParseIP("::1"),
		},
		ipv6Err: nil,
	}
	selectionOption := "all"
	useIPv4 := true

	// act
	ips, _ := getMyIP(ipProvider, selectionOption, useIPv4)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("127.0.0.1"),
	}
	if fmt.Sprintf("%s", ips) != fmt.Sprintf("%s", expectedResult) {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) returned %q but should have returned %q", selectionOption, useIPv4, ips, expectedResult)
	}

}

// getMyIP should return the IPv6 addresses of the IP provider if the useIPv4 flag is set to false.
func Test_getMyIP_UseIPv4False_IPProviderHasIPv6Addresses_IPv46AddressesAreReturned(t *testing.T) {
	// arrange
	ipProvider := testIPProvider{
		ipv4IPs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		ipv4Err: nil,
		ipv6IPs: []net.IP{
			net.ParseIP("::1"),
		},
		ipv6Err: nil,
	}
	selectionOption := "all"
	useIPv4 := false

	// act
	ips, _ := getMyIP(ipProvider, selectionOption, useIPv4)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("::1"),
	}
	if fmt.Sprintf("%s", ips) != fmt.Sprintf("%s", expectedResult) {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) returned %q but should have returned %q", selectionOption, useIPv4, ips, expectedResult)
	}

}

// getMyIP should not return anything if the IP provider has no IPv6 addresses.
func Test_getMyIP_UseIPv4False_IPProviderHasNoIPv6Addresses_ResultIsEmpty(t *testing.T) {
	// arrange
	ipProvider := testIPProvider{
		ipv4IPs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		ipv4Err: nil,
	}
	selectionOption := "all"
	useIPv4 := false

	// act
	ips, _ := getMyIP(ipProvider, selectionOption, useIPv4)

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) returned %q but should not have returned anything because the IP provider has no IPv6 addresses", selectionOption, useIPv4, ips)
	}

}

// getMyIP (useIPv4 = false) should only return an error if the IP provider returns an error.
func Test_getMyIP_UseIPv4False_IPProviderReturnsError_NoIPsAreReturned_ErrorIsReturned(t *testing.T) {
	// arrange
	ipProvider := testIPProvider{
		ipv4IPs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		ipv4Err: fmt.Errorf("IPv4 error"),
		ipv6IPs: []net.IP{
			net.ParseIP("::1"),
		},
		ipv6Err: fmt.Errorf("IPv6 error"),
	}
	selectionOption := "all"
	useIPv4 := false

	// act
	ips, err := getMyIP(ipProvider, selectionOption, useIPv4)

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) returned %q but should not have returned anything because the IP provider returned an error", selectionOption, useIPv4, ips)
	}

	if err == nil {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) did not return an error even though the IP Provider responded with one.", selectionOption, useIPv4)
	}

}

// getMyIP (useIPv4 = true) should only return an error if the IP provider returns an error.
func Test_getMyIP_UseIPv4True_IPProviderReturnsError_NoIPsAreReturned_ErrorIsReturned(t *testing.T) {
	// arrange
	ipProvider := testIPProvider{
		ipv4IPs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		ipv4Err: fmt.Errorf("IPv4 error"),
		ipv6IPs: []net.IP{
			net.ParseIP("::1"),
		},
		ipv6Err: fmt.Errorf("IPv6 error"),
	}
	selectionOption := "all"
	useIPv4 := false

	// act
	ips, err := getMyIP(ipProvider, selectionOption, useIPv4)

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) returned %q but should not have returned anything because the IP provider returned an error", selectionOption, useIPv4, ips)
	}

	if err == nil {
		t.Fail()
		t.Logf("getMyIP(ipProvider, %q, %v) did not return an error even though the IP Provider responded with one.", selectionOption, useIPv4)
	}

}

// If no IPs are supplied and no select option no error should be returned.
func Test_getSelectedIPs_NoIPsSupplied_NoSelectOptionSupplied_ResultIsEmpty_NoError(t *testing.T) {
	// arrange
	ips := []net.IP{}
	selectOption := ""

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	if len(selectedIPs) > 0 {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q. But the result should be empty.", ips, selectOption, selectedIPs)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not have returned %q.", ips, selectOption, err.Error())
	}
}

// If no IPs are supplied but a select option is given an error should be returned.
func Test_getSelectedIPs_NoIPsSupplied_SelectOptionSupplied_ErrorIsReturned(t *testing.T) {
	// arrange
	ips := []net.IP{}
	selectOption := "all"

	// act
	_, err := getSelectedIPs(ips, selectOption)

	// assert
	if err == nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should return an error.", ips, selectOption)
	}
}

// Invalid select options should result in an error.
func Test_getSelectedIPs_SelectOptionIsInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	ips := []net.IP{}
	invalidOptions := []string{
		"  ",
		" all",
		"all ",
		"dasdsadsa",
		"1 2 3",
		"1;2;3",
		"1,",
		",1,2,3",
	}

	for _, selectOption := range invalidOptions {

		// act
		_, err := getSelectedIPs(ips, selectOption)

		// assert
		if err == nil {
			t.Fail()
			t.Errorf("getSelectedIPs(%q, %q) should return an error because the given option is invalid.", ips, selectOption)
		}
	}
}

// If the select option "all" is used all IPs should be returned.
func Test_getSelectedIPs_IPsSupplied_SelectOptionAll_AllIPsAreReturned(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "all"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	if fmt.Sprintf("%s", ips) != fmt.Sprintf("%s", selectedIPs) {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, ips)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "first" is used only the first IP should be returned.
func Test_getSelectedIPs_IPsSupplied_SelectOptionFirst_FirstIPIsReturned(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "first"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	if len(selectedIPs) != 1 || fmt.Sprintf("%s", selectedIPs[0]) != "127.0.0.1" {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, ips[:1])
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "last" is used only the last IP should be returned.
func Test_getSelectedIPs_IPsSupplied_SelectOptionLast_LastIPIsReturned(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "last"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	if len(selectedIPs) != 1 || fmt.Sprintf("%s", selectedIPs[0]) != "127.0.0.3" {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, ips[2:3])
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "1,2,3" is used the IPs should be returned in the specified order.
func Test_getSelectedIPs_IPsSupplied_SelectOption123_IPsAreReturnedInCorrectOrder(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "1,2,3"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}

	if fmt.Sprintf("%s", expectedResult) != fmt.Sprintf("%s", selectedIPs) {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, expectedResult)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "3,2,1" is used the IPs should be returned in the specified order.
func Test_getSelectedIPs_IPsSupplied_SelectOption321_IPsAreReturnedInCorrectOrder(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "3,2,1"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("127.0.0.3"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.1"),
	}

	if fmt.Sprintf("%s", expectedResult) != fmt.Sprintf("%s", selectedIPs) {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, expectedResult)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "3,1" is used the IPs should be returned in the specified order.
func Test_getSelectedIPs_IPsSupplied_SelectOption31_IPsAreReturnedInCorrectOrder(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "3,1"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("127.0.0.3"),
		net.ParseIP("127.0.0.1"),
	}

	if fmt.Sprintf("%s", expectedResult) != fmt.Sprintf("%s", selectedIPs) {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, expectedResult)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option "1,1,1,1" is used the IPs should be returned in the specified order.
// Returning the same IP multiple times should be possible (even though I don't know why you would want that).
func Test_getSelectedIPs_IPsSupplied_SelectOption1111_IPsAreReturnedInCorrectOrder(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := "1,1,1,1"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	expectedResult := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.1"),
	}

	if fmt.Sprintf("%s", expectedResult) != fmt.Sprintf("%s", selectedIPs) {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should have returned %q.", ips, selectOption, selectedIPs, expectedResult)
	}

	if err != nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should not return an error but returned: %s.", ips, selectOption, err.Error())
	}
}

// If the select option ",1" is used no IPs should be returned but an error.
func Test_getSelectedIPs_IPsSupplied_SelectOptionIsInvalid_NoIPsAreReturned_ErrorIsReturned(t *testing.T) {
	// arrange
	ips := []net.IP{
		net.ParseIP("127.0.0.1"),
		net.ParseIP("127.0.0.2"),
		net.ParseIP("127.0.0.3"),
	}
	selectOption := ",1"

	// act
	selectedIPs, err := getSelectedIPs(ips, selectOption)

	// assert
	if len(selectedIPs) > 0 {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) returned %q but should not have returned any IPs because the select option is invalid.", ips, selectOption, selectedIPs)
	}

	if err == nil {
		t.Fail()
		t.Errorf("getSelectedIPs(%q, %q) should return an error but did not.", ips, selectOption)
	}
}
