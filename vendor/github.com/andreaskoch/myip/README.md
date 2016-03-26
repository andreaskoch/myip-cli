# myip

A small go package that returns your current local or remote IP addresses

**myip** is a small go package that just **returns your current IP addresses** (local or remote).

## Usage

Use `github.com/andreaskoch/myip` in your program.

### Get your local IP addresses

```go
package main

import (
	"fmt"
	"github.com/andreaskoch/myip"
	"os"
)

func main() {

	// create a new local IP provider instance
	remoteIPProvider, ipProviderError := myip.NewLocalIPProvider()
	if ipProviderError != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a new local IP provider: %s", ipProviderError.Error())
		os.Exit(1)
	}

	// get the local IPv6 addresses
	remoteIPv6Addresses, remoteIPError := remoteIPProvider.GetIPv6Addresses()
	if remoteIPError != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine the local IPv6 addresses: %s", remoteIPError.Error())
		os.Exit(1)
	}

	// print the the local IPv6 addresses
	fmt.Fprintf(os.Stdout, "%s", remoteIPv6Addresses)

}
```

### Get your remote IP address

```go
package main

import (
	"fmt"
	"github.com/andreaskoch/myip"
	"os"
)

func main() {

	// create a new remote IP provider instance
	remoteIPProvider, ipProviderError := myip.NewRemoteIPProvider()
	if ipProviderError != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a new remote IP provider: %s", ipProviderError.Error())
		os.Exit(1)
	}

	// get the remote IPv6 addresses
	remoteIPv6Addresses, remoteIPError := remoteIPProvider.GetIPv6Addresses()
	if remoteIPError != nil {
		fmt.Fprintf(os.Stderr, "Failed to determine the remote IPv6 addresses: %s", remoteIPError.Error())
		os.Exit(1)
	}

	// print the the remote IPv6 addresses
	fmt.Fprintf(os.Stdout, "%s", remoteIPv6Addresses)

}
```

## Remote IP services

For determining your remote IP address (IPV6 or IPv4) myip relies on external services:

1. yip.li
	- IPv6: [ipv6.yip.li](https://ipv6.yip.li)
	- IPv4: [ipv4.yip.li](https://ipv4.yip.li)
2. icanhazip.com
	- IPv6: [ipv6.icanhazip.com](http://ipv6.icanhazip.com)
	- IPv4: [ipv4.icanhazip.com](http://ipv4.icanhazip.com)

**yip.li** is developed and hosted by me ([Andreas Koch](https://andykdocs.de/about)). You can find the source code at [github.com/andreaskoch/yip](https://github.com/andreaskoch/yip). Information about **icanhazip.com** can be found at [github.com/major/icanhaz](https://github.com/major/icanhaz).

For determining the IP address myip will call both services and whoever responds first will provide the your remote IP.

## Roadmap

### Trusted Remote IP Detection

Make sure the provided remote IPs are correct and cannot be manipulated by manipulating one of the remote IP services.

- SSL certificate pinning for the remote services
	- Maybe. I am not sure if this is worth the effort.

and/or

- Add more remote services and introduce a majority decision.
	- If the majority of the services return the same IP that one will be taken.
	- If there is no majority and error is reported.
