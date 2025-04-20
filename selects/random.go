package selects

import (
	"fmt"
	"math/rand/v2"

	"github.com/nezbut/proxym"
)

// RandomSelect is a proxy selection strategy that returns a random proxy.
type RandomSelect struct {
	provider proxym.SelectStrategyProxyProvider
}

// NewRandomSelect returns a new RandomSelect.
func NewRandomSelect(provider proxym.SelectStrategyProxyProvider) proxym.SelectStrategy {
	return &RandomSelect{
		provider: provider,
	}
}

// Select returns the proxy to use.
func (s *RandomSelect) Select() (*proxym.Proxy, error) {
	proxies := s.provider.GetProxies()
	if len(proxies) == 0 {
		return nil, fmt.Errorf("%w: empty proxies from provider", proxym.ErrFailedSelectProxy)
	}
	return proxies[rand.IntN(len(proxies))], nil //nolint: gosec // can be used ordinary random sampling
}
