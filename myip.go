// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// GitInfo is either the empty string (the default)
// or is set to the git hash of the most recent commit
// using the -X linker flag (Example: "2015-01-11-284c030+")
var GitInfo string

// commandOptions is the flag set for the "local" and "remote" actions
var commandOptions = flag.NewFlagSet("command-options", flag.ExitOnError)

// useIPv4 contains a flag inidicating whether IPv4 addresses should be used (default: false)
var useIPv4 bool

// ipSelectionOption specifies the IP address that shall be returned if there are multiple addresses available
var ipSelectionOption string

// ipSelectionOptionIndexPattern defines the pattern for the index-based IP selection (e.g. "1,2,3")
var ipSelectionOptionIndexPattern = regexp.MustCompile(`^(\d+,*)+$`)

const ipSelectionOptionAll = "all"
const ipSelectionOptionFirst = "first"
const ipSelectionOptionLast = "last"

var ipSelectionOptions = []string{ipSelectionOptionAll, ipSelectionOptionFirst, ipSelectionOptionLast, "1", "1,3"}

// actionnamelocal contains the name of the "local" action
const actionnamelocal = "local"

// actionnameremote contains the name of the "remote" action
const actionnameremote = "remote"

// The ipAddresser interface provides functions for
// retrieving IPv4 and IPv6 addresses.
type ipAddresser interface {
	ipv4Addresser
	ipv6Addresser
}

// The ipv6Addresser interface provides functions for
// retrieving IPv6 addresses.
type ipv6Addresser interface {
	GetIPv6Addresses() ([]net.IP, error)
}

// The ipv4Addresser interface provides functions for
// retrieving IPv4 addresses.
type ipv4Addresser interface {
	GetIPv4Addresses() ([]net.IP, error)
}

func main() {

	arguments := os.Args
	executableName := arguments[0]

	commandOptions.BoolVar(&useIPv4, "4", false, fmt.Sprintf("Use IPv4 instead of IPv6"))
	commandOptions.StringVar(&ipSelectionOption, "select", ipSelectionOptionAll, fmt.Sprintf("Select one or more IPs (\"%s\")", strings.Join(ipSelectionOptions, `", "`)))

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s returns your local IPv6 (or IPv4) address.\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Version: %s\n", version())
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  %s <action> [options]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Actions:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "%10s  %s\n", actionnamelocal, "Get your local IP address")
		fmt.Fprintf(os.Stderr, "%10s  %s\n", actionnameremote, "Get your remote IP address")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		commandOptions.PrintDefaults()
	}

	flag.Parse()

	// get action
	if len(arguments) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// parse the command line options
	commandOptions.Parse(arguments[2:])

	// action: remote vs. local
	var ips []net.IP
	var myIPError error

	actionName := strings.TrimSpace(strings.ToLower(arguments[1]))
	switch actionName {
	case actionnamelocal:
		ips, myIPError = myLocalIP(ipSelectionOption, useIPv4)

	case actionnameremote:
		ips, myIPError = myRemoteIP(ipSelectionOption, useIPv4)

	default:
		{
			fmt.Fprintf(os.Stderr, "The action %q does not exist.\n\n", actionName)
			flag.Usage()
			os.Exit(1)
		}
	}

	// print errors
	if myIPError != nil {
		fmt.Fprintf(os.Stderr, "%s", myIPError.Error())
		os.Exit(1)
	}

	// print IPs
	for _, ip := range ips {
		fmt.Fprintf(os.Stdout, "%s\n", ip)
	}

}

// myLocalIP returns the current local IPv6 (or IPv4) address
func myLocalIP(selectionOption string, useIPv4 bool) ([]net.IP, error) {

	ipProvider, ipProviderError := newLocalIPProvider()
	if ipProviderError != nil {
		return nil, fmt.Errorf("%s\n", ipProviderError.Error())
	}

	return getMyIP(ipProvider, selectionOption, useIPv4)
}

// myRemoteIP returns the current remote IPv6 (or IPv4) address
func myRemoteIP(selectionOption string, useIPv4 bool) ([]net.IP, error) {

	ipProvider, ipProviderError := newRemoteIPProvider()
	if ipProviderError != nil {
		return nil, fmt.Errorf("%s\n", ipProviderError.Error())
	}

	return getMyIP(ipProvider, selectionOption, useIPv4)
}

// getMyIP returns the selected IPv6 or IPv4 addresses from the given IP provider.
func getMyIP(ipProvider ipAddresser, selectionOption string, useIPv4 bool) ([]net.IP, error) {

	// IPv6 vs IPv4
	var allIPs []net.IP
	var ipErr error
	if useIPv4 {
		allIPs, ipErr = ipProvider.GetIPv4Addresses()
	} else {
		allIPs, ipErr = ipProvider.GetIPv6Addresses()
	}

	// handle errors
	if ipErr != nil {
		return nil, fmt.Errorf("%s\n", ipErr.Error())
	}

	// abort if no IPs are returned
	if len(allIPs) == 0 {
		ipType := "IPv6"
		if useIPv4 {
			ipType = "IPv4"
		}

		return []net.IP{}, fmt.Errorf("No %s IPs available.", ipType)
	}

	// select one or more IPs
	selectedIPs, ipSelectionError := getSelectedIPs(allIPs, selectionOption)
	if ipSelectionError != nil {
		return nil, fmt.Errorf("%s\n", ipSelectionError.Error())
	}

	return selectedIPs, nil
}

// getSelectedIPs returns a subset of the given IPs based on the given selection option (all, fist, last, "1,2", ...).
// If the given selection option is invalid an error will be returned.
func getSelectedIPs(ips []net.IP, selectionOption string) ([]net.IP, error) {

	// abort if no IPs have been supplied
	if len(ips) == 0 {

		// If there was a selection given, an empty IP slice is an error
		selectionGiven := len(selectionOption) > 0
		selectOptionIsNotAll := selectionOption != ipSelectionOptionAll
		if selectionGiven && selectOptionIsNotAll {
			return []net.IP{}, fmt.Errorf("Invalid selection %q. No IPs available.", selectionOption)
		}

		// no error
		return []net.IP{}, nil

	}

	// handle special options: all, first, last
	switch {
	case selectionOption == ipSelectionOptionAll:
		return ips, nil

	case selectionOption == ipSelectionOptionFirst:
		return ips[:1], nil

	case selectionOption == ipSelectionOptionLast:
		{
			lastIPIndex := len(ips) - 1
			return ips[lastIPIndex:], nil
		}

	case ipSelectionOptionIndexPattern.MatchString(selectionOption):
		{
			break
		}

	default:
		{
			return []net.IP{}, fmt.Errorf("%q is not a valid value for the IP selection", ipSelectionOption)
		}
	}

	// handle indexed selection
	var selectedIPs []net.IP
	selectedIndizes := strings.Split(selectionOption, ",")
	for _, indexString := range selectedIndizes {

		// parse the string
		index64, err := strconv.ParseInt(indexString, 10, 64)
		if err != nil {
			return []net.IP{}, fmt.Errorf("Invalid IP selection index supplied (min: 1, max: %d).\n", len(ips))
		}

		index := int(index64)

		// verify the index
		if index < 1 || index > len(ips) {
			return []net.IP{}, fmt.Errorf("Invalid IP selection index supplied (min: 1, max: %d).\n", len(ips))
		}

		// append the selected IP
		selectedIPs = append(selectedIPs, ips[index-1])

	}

	return selectedIPs, nil
}

// version returns the git version of this binary (e.g. "2015-01-11-284c030+").
// If the linker flags were not provided, the return value is "unknown".
func version() string {
	if GitInfo != "" {
		return GitInfo
	}

	return "unknown"
}
