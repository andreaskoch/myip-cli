// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"testing"
)

func Test_GetIPs_NoInterfacesSupplied_NoIPsAreReturned(t *testing.T) {
	// arrange
	interfaces := []net.Interface{}
	addressProvider := interfaceAddressProvider{interfaces}

	// act
	ips, _ := addressProvider.GetIPs()

	// assert
	if len(ips) > 0 {
		t.Fail()
		t.Logf("GetIPs() should return an empty list of IPs but returned %v instead.", ips)
	}
}
