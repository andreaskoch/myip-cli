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
