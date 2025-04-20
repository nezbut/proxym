package proxym

// SelectStrategy is an interface for proxy selection strategies.
// It is used to determine which proxy to use.
type SelectStrategy interface {
	// Select returns the proxy to use.
	//
	// If the strategy fails to select a proxy, an ErrFailedSelectProxy error is returned.
	Select() (*Proxy, error)
}

// SelectStrategyProxyProvider is an interface for proxy selection strategies providers.
//
// Used to get a list of proxies to choose from.
type SelectStrategyProxyProvider interface {
	// GetProxies returns the copied list of proxies.
	GetProxies() []*Proxy
}

// SelectStrategyFactory is a function that returns a SelectStrategy from a SelectStrategyProxyProvider.
type SelectStrategyFactory func(SelectStrategyProxyProvider) SelectStrategy
