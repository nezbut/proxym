# Proxym

[![Go Reference](https://pkg.go.dev/badge/github.com/nezbut/proxym.svg)](https://pkg.go.dev/github.com/nezbut/proxym)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![CI Status](https://github.com/nezbut/proxym/actions/workflows/ci.yml/badge.svg)](https://github.com/nezbut/proxym/actions)

**Proxym** is a library that provides a flexible and customizable proxy manager for Go applications. It simplifies proxy rotation and selection with a system of strategies and filters, providing thread safety and seamless integration with standard `http.RoundTripper` interfaces.

Proxym provides the main interface: `proxym.ProxyManager`, which allows you to manage and select proxies efficiently.

Also in proxym there is `proxym.ProxyManagerImpl` which you create via `proxym.NewProxyManager(...)` this component is a standard implementation for ProxyManager which uses SelectStrategy and RotationStrategy to get the next available proxy and also there you can configure resource-speciffic proxies and strategies

You can implement your own custom ProxyManager which should correspond to the `proxym.ProxyManager` interface and for example pass it to `proxym.NewClient(customProxyManager)`.

## Features

- **Global and resource-specific configurations**: define proxies and strategies at global proxy manager and per-resource levels.
- **Proxy statistics and metadata**: view proxy statistics and manage metadata.
- **Managed proxies**: disable/enable, manage metadata, view is active/direct.
- **Select strategies**: determine which proxy to use.
- **Rotation strategies**: determine if a proxy should be rotated.
- **Select filters**: apply filters before selection.
- **HTTP integration**: use with any HTTP client that supports `http.RoundTripper`.
- **Thread-safe**: thread-safe for concurrent use.

## Installation

```bash
go get github.com/nezbut/proxym
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/nezbut/proxym"
	"github.com/nezbut/proxym/rotations"
	"github.com/nezbut/proxym/selects"
)

func main() {
	// Define a list of proxies
	proxies := []*proxym.Proxy{
		proxym.NewProxyStr("http://proxy1:8080", nil),
		proxym.NewProxyStr("socks5://localhost:9050", nil),
		// direct connection
		proxym.NewDirectConnection(), // use this if you want to use your direct connection during requests.
	}

	// Create a proxy manager
	pm := proxym.NewProxyManager(
		// proxym.WithProxies(proxies...), add proxies in options
		proxym.WithSelectStrategy(selects.DefaultSelectStrategy()),       // default select strategy
		proxym.WithRotationStrategy(rotations.DefaultRotationStrategy()), // default rotation strategy
	)

	pm.AddProxies(proxies...) // add proxies in runtime

	// Create a http client
	client := proxym.NewClient(pm) // or use proxym.PatchClient(client, pm) to patch existing http client

	// Perform requests
	req, _ := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://api.ipify.org/",
		nil,
	)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Last used proxy:", pm.LastUsed())
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Response status code:", resp.StatusCode)
	log.Println("Last used proxy:", pm.LastUsed())
	log.Println("Is direct connection:", pm.LastUsed().IsDirect())
	// After startup you may get 3 responses
	//
	// 1. This is an error and Last used proxy will be http://proxy1:8080 unless of course you have one
	//
	// 2. There will be a 200 response and Last used proxy will be socks5://localhost:9050
	// of course if you have tor running and enabled, direct connection will be false
	//
	// 3. The answer will be 200 and Last used proxy will be <not proxy url> direct connection will be true
}

```

## Strategies

Strategies for `proxym.ProxyManagerImpl`

### Rotation

Strategies that determine whether to rotate proxies, package: `proxym/rotations`

#### Strategies realizations

- `rotations.CompositeRotation`: combines multiple rotation strategies with a specified logic(AND or OR).
- `rotations.OnlyEnabledRotation`: returns true if the proxy is disabled.
- `rotations.ErrorThresholdRotation`: returns true if the error proxy is greater than or equal to a threshold.
- `rotations.RequestLimitedRotation`: returns true if the total number of requests is greater than or equal to a limit.
- `rotations.RoundRobinRotation`: always returns true.

Default rotation strategy get from `rotations.DefaultRotationStrategy()`

For create custom rotation strategy implement the `proxym.RotationStrategy` interface.

### Select

Strategies that determine which proxy to use, package: `proxym/selects`

#### Strategies realizations

- `selects.RoundRobinSelect`: returns proxies in a round-robin fashion.
- `selects.RandomSelect`: returns a random proxy.

Default select strategy get from `selects.DefaultSelectStrategy()`

For create custom select strategy implement the `proxym.SelectStrategy` interface and create `proxym.SelectStrategyFactory` for this implementation.

To pass your select strategy in options in `proxym.ProxyManager` or `proxym.ResourceConfig`,
you need to create a factory function that will return `proxym.SelectStrategy`,
for example `selects.RandomSelect` is `selects.NewRandomSelect`, this function simply takes `proxym.SelectStrategyProxyProvider` and creates a new `selects.RandomSelect` strategy

`proxym.SelectStrategyProxyProvider` is an interface that provides a proxy from somewhere for the select strategy.
For example, its implementations are `proxym.ProxyManager` and `proxym.ResourceConfig`

### Select filters

Filters that modify the list of proxies before selection, package: `proxym/selects`

#### Filters realizations

- `selects.RemoveDisabledFilter`: excludes proxies marked as disabled.
- `selects.RemoveActiveProxyFilter`: excludes the active proxy to avoid repetition.

For create custom select filter implement the `selects.SelectFilter` interface.

Example of how to create SelectStrategy with filters

```go
package main

import (
	"github.com/nezbut/proxym"
	"github.com/nezbut/proxym/selects"
)

// Returns a RandomSelect with RemoveActiveProxyFilter and RemoveDisabledFilter filters.
// Function returns proxym.SelectStrategyFactory.
func getSelectStrategy() proxym.SelectStrategyFactory {
	return selects.NewFilteredSelectFactory(
		// The factory of select strategies used.
		selects.NewRandomSelect, // Here the factory for selects.RandomSelect is used, you can also use selects.NewRoundRobinSelect.
		// You can pass any number of filters here.
		selects.RemoveActiveProxyFilter{},
		selects.RemoveDisabledFilter{},
		// ...
	)
}

func main() {

	// Create a resource config with select strategy and filters
	proxym.NewResourceConfig(
		// ...
		proxym.WithResourceSelectStrategy(getSelectStrategy()),
		// ...
	)

	// Create a proxy manager with select strategy and filters
	proxym.NewProxyManager(
		// ...
		proxym.WithSelectStrategy(getSelectStrategy()),
		// ...
	)

	// ...
}

```

## Advanced Usage

### Resource-specific proxies and strategies

You can configure different proxies and strategies for different resources.

```go
package main

import (
	"github.com/nezbut/proxym"
	"github.com/nezbut/proxym/rotations"
	"github.com/nezbut/proxym/selects"
)

func main() {
	// Define a list of proxies
	proxies := []*proxym.Proxy{
		proxym.NewProxyStr("http://proxy1:8080", nil),
		proxym.NewProxyStr("socks5://localhost:9050", nil),
	}

	// Create a resource config for https://ipify.org/
	resource := proxym.NewResourceConfig(
		true, // normalize domain, https://ipify.org/ will be normalized to ipify.org
		proxym.WithDomain("https://ipify.org/"),
		// proxym.WithResourceProxies(proxies...), add proxies in options
		proxym.WithResourceSelectStrategy(selects.DefaultSelectStrategy()),
		proxym.WithResourceRotationStrategy(rotations.DefaultRotationStrategy()),
		// ignore subdomains in the comparison of the domain
		// proxym.WithIgnoreSubdomains(false), // default true, if is false, then api.ipify.org != ipify.org
	)
	resource.AddProxies(proxies...) // add proxies in runtime

	// Create a proxy manager
	pm := proxym.NewProxyManager(
		proxym.WithResources(resource),                                   // add resource
		proxym.WithSelectStrategy(selects.DefaultSelectStrategy()),       // default select strategy
		proxym.WithRotationStrategy(rotations.DefaultRotationStrategy()), // default rotation strategy
	)

	// pm.AddResourceProxies("ipify.org", proxies...) // add proxies in runtime by proxy manager

	// Create a http client
	client := proxym.NewClient(pm) // or use proxym.PatchClient(client, pm) to patch existing http client

	// Perform requests...
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
