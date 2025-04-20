package proxym

import (
	"errors"
	"fmt"
	"sync"
)

// ProxyManager is a manager for proxies.
//
// Allows you to receive the following proxy depending on the internal state.
// You can get the last used proxy and the list of proxies.
type ProxyManager interface {
	// GetNextProxy returns the next available proxy by domain.
	GetNextProxy(domain string) (*Proxy, error)
	// LastUsed Returns the last used proxy.
	// This method may return nil in *Proxy if no proxy has been used.
	LastUsed() *Proxy
	// GetProxies returns the copied list of proxies.
	GetProxies() []*Proxy
}

// ProxyManagerImpl is a ProxyManager implementation.
type ProxyManagerImpl struct {
	proxies          []*Proxy
	pMu              sync.RWMutex
	resources        []*ResourceConfig
	rMu              sync.RWMutex
	lastUsed         *Proxy
	rotationStrategy RotationStrategy
	selectStrategy   SelectStrategy
	mu               sync.RWMutex
}

// NewProxyManager creates a new ProxyManagerImpl.
//
// Important:
//   - Proxy manager starts with empty proxy list and nothing strategies by default
//   - You MUST add proxies using either:
//   - WithProxies() option during initialization
//   - ProxyManagerImpl.AddProxies()
//   - You MUST add strategies using either:
//   - WithRotationStrategy() option during initialization
//   - WithSelectStrategy() option during initialization
//   - If you don't set strategies, the constructor will panic
//
// Example minimum working setup:
//
//	proxy2, _ := url.Parse("http://proxy2:8080")
//	proxy3, _ := proxym.NewProxyParsedStr("http://proxy3:8080", nil)
//
//	pm := proxym.NewProxyManager(
//	    proxym.WithProxies(
//	        proxym.NewProxyStr("http://proxy1:8080", nil),
//	        proxym.NewProxy(proxy2, nil),
//	        proxy3,
//	        proxym.NewDirectConnection(), // is not proxy, is direct connection.
//	    ),
//	    proxym.WithRotationStrategy(rotations.DefaultRotationStrategy()),
//	    proxym.WithSelectStrategy(selects.DefaultSelectStrategy()),
//	)
func NewProxyManager(opts ...ProxyManagerImplOption) *ProxyManagerImpl {
	pm := &ProxyManagerImpl{
		proxies:   make([]*Proxy, 0),
		resources: make([]*ResourceConfig, 0),
	}
	for _, opt := range opts {
		opt(pm)
	}
	if pm.rotationStrategy == nil || pm.selectStrategy == nil {
		panic("rotationStrategy and selectStrategy must be set")
	}
	return pm
}

// GetNextProxy returns the next available proxy.
// If the resource by domain is not found global is returned.
//
// If SelectStrategy returns nil and err is nil, then there will be an error ErrProxyNotAvailable.
func (pm *ProxyManagerImpl) GetNextProxy(domain string) (*Proxy, error) {
	if len(pm.proxies) == 0 && len(pm.resources) == 0 {
		return nil, pm.proxyNotAvailable(ErrEmptyProxyList)
	}
	resource, err := pm.getResourceByDomain(domain)
	isNotFound := errors.Is(err, ErrResourceNotFound)
	if err != nil && !isNotFound {
		return nil, pm.proxyNotAvailable(err)
	}
	lastUsed := pm.LastUsed()
	var current *Proxy

	if isNotFound { //nolint:nestif // don't
		if lastUsed != nil && !pm.rotationStrategy.ShouldRotate(lastUsed) {
			return lastUsed, nil
		}

		currentProxy, errSelect := pm.selectStrategy.Select()
		if errSelect != nil {
			return nil, pm.proxyNotAvailable(errSelect)
		}

		current = currentProxy
	} else {
		if lastUsed != nil && !resource.rotationStrategy.ShouldRotate(lastUsed) {
			return lastUsed, nil
		}

		currentProxy, errSelect := resource.selectStrategy.Select()
		if errSelect != nil {
			return nil, pm.proxyNotAvailable(errSelect)
		}

		current = currentProxy
	}

	if current == nil {
		return nil, ErrProxyNotAvailable
	}

	if lastUsed != nil {
		lastUsed.deactivate()
	}
	current.activate()

	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.lastUsed = current
	return current, nil
}

// LastUsed Returns the last used proxy.
// This method may return nil in *Proxy if no proxy has been used.
func (pm *ProxyManagerImpl) LastUsed() *Proxy {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.lastUsed
}

// GetProxies returns the copied list of proxies.
func (pm *ProxyManagerImpl) GetProxies() []*Proxy {
	pm.pMu.RLock()
	defer pm.pMu.RUnlock()

	proxies := make([]*Proxy, len(pm.proxies))
	copy(proxies, pm.proxies)

	return proxies
}

// AddResources adds resources to the ProxyManagerImpl.
func (pm *ProxyManagerImpl) AddResources(resources ...*ResourceConfig) {
	pm.rMu.Lock()
	defer pm.rMu.Unlock()
	pm.resources = append(pm.resources, resources...)
}

// AddProxies adds proxies to the ProxyManagerImpl.
func (pm *ProxyManagerImpl) AddProxies(proxies ...*Proxy) {
	pm.pMu.Lock()
	defer pm.pMu.Unlock()
	pm.proxies = append(pm.proxies, proxies...)
}

// AddResourceProxies adds proxies to the ResourceConfig by domain.
func (pm *ProxyManagerImpl) AddResourceProxies(domain string, proxies ...*Proxy) error {
	resource, err := pm.getResourceByDomain(domain)

	if err != nil {
		return err
	}

	resource.AddProxies(proxies...)
	return nil
}

func (pm *ProxyManagerImpl) getResourceByDomain(domain string) (*ResourceConfig, error) {
	pm.rMu.RLock()
	defer pm.rMu.RUnlock()

	for _, resource := range pm.resources {
		if resource.CompareDomain(domain) {
			return resource, nil
		}
	}
	return nil, ErrResourceNotFound
}

func (pm *ProxyManagerImpl) proxyNotAvailable(err error) error {
	return fmt.Errorf("%w: %w", ErrProxyNotAvailable, err)
}
