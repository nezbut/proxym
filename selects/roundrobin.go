package selects

import (
	"fmt"
	"sync"

	"github.com/nezbut/proxym"
)

// RoundRobinSelect is a proxy selection strategy that returns proxies in a round-robin fashion.
//
// The first time Select is called, it will return the first proxy from the provider.
// Each subsequent call to Select will return the next proxy from the provider
// until the end of the list is reached, at which point it will start from the beginning again.
type RoundRobinSelect struct {
	provider proxym.SelectStrategyProxyProvider
	index    int
	mu       sync.Mutex
}

// NewRoundRobinSelect returns a new RoundRobinSelect.
//
// The index is set to -1, so the first call to Select() will start with the first proxy.
func NewRoundRobinSelect(provider proxym.SelectStrategyProxyProvider) proxym.SelectStrategy {
	return &RoundRobinSelect{
		provider: provider,
		index:    -1,
	}
}

// Select returns the proxy to use.
func (s *RoundRobinSelect) Select() (*proxym.Proxy, error) {
	proxies := s.provider.GetProxies()
	if len(proxies) == 0 {
		return nil, fmt.Errorf("%w: empty proxies from provider", proxym.ErrFailedSelectProxy)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.index = (s.index + 1) % len(proxies)
	return proxies[s.index], nil
}
