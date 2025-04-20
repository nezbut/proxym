package proxym

import (
	"net/url"
	"strings"
	"sync"
)

// ResourceConfig is a representation of a resource config in proxym.
//
// These are the proxy, RotationStrategy and SelectStrategy settings for a particular resource.
type ResourceConfig struct {
	proxies             []*Proxy
	domain              string
	notIgnoreSubdomains bool
	selectStrategy      SelectStrategy
	rotationStrategy    RotationStrategy
	mu                  sync.RWMutex
}

// NewResourceConfig creates a new ResourceConfig.
//
// If normalizeDomain is true, the domain will be normalized.
//
// Important:
//   - Resource config starts with empty proxy list and nothing strategies by default
//   - You MUST add proxies using either:
//   - WithResourceProxies() option during initialization
//   - ProxyManagerImpl.AddResourceProxies()
//   - You MUST add strategies using either:
//   - WithResourceRotationStrategy() option during initialization
//   - WithResourceSelectStrategy() option during initialization
//   - If you don't set strategies, the constructor will panic
//
// Example minimum working setup:
//
//	proxy2, _ := url.Parse("http://proxy2:8080")
//
//	rc := proxym.NewResourceConfig(
//	    true, // normalize domain, http://api.example.com will be normalized to api.example.com
//	    proxym.WithDomain("api.example.com"),
//	    proxym.WithResourceRotationStrategy(rotations.DefaultRotationStrategy()),
//	    proxym.WithResourceSelectStrategy(selects.DefaultSelectStrategy()),
//	    proxym.WithResourceProxies(
//	        proxym.NewProxyStr("http://proxy1:8080", nil),
//	        proxym.NewProxy(proxy2, nil),
//	    ),
//	)
func NewResourceConfig(normalizeDomain bool, opts ...ResourceConfigOption) *ResourceConfig {
	rc := &ResourceConfig{
		proxies: make([]*Proxy, 0),
	}

	for _, opt := range opts {
		opt(rc)
	}

	if rc.rotationStrategy == nil || rc.selectStrategy == nil {
		panic("RotationStrategy and SelectStrategy must be set")
	}

	if normalizeDomain {
		rc.domain = rc.normalizeDomain(rc.domain)
	}
	return rc
}

// Domain returns the domain of the ResourceConfig.
func (rc *ResourceConfig) Domain() string {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.domain
}

// GetProxies returns the copied list of proxies.
func (rc *ResourceConfig) GetProxies() []*Proxy {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	proxies := make([]*Proxy, len(rc.proxies))
	copy(proxies, rc.proxies)

	return proxies
}

// AddProxies adds proxies to the ResourceConfig.
func (rc *ResourceConfig) AddProxies(proxies ...*Proxy) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.proxies = append(rc.proxies, proxies...)
}

// CompareDomain compare domain.
//
// If notIgnoreSubdomains is false, then it will ignore subdomains in the comparison of the domain.
func (rc *ResourceConfig) CompareDomain(domain string) bool {
	rcDomain := rc.Domain()
	normalized := rc.normalizeDomain(domain)

	if normalized == rcDomain {
		return true
	}

	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if !rc.notIgnoreSubdomains && strings.HasSuffix(normalized, "."+rcDomain) {
		return true
	}

	return false
}

// normalizeDomain normalizes domain.
func (rc *ResourceConfig) normalizeDomain(domain string) string {
	if domain == "" {
		return ""
	}
	return strings.ToLower(rc.getDomainFromURL(domain))
}

// getDomainFromURL gets domain from url.
func (rc *ResourceConfig) getDomainFromURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil || u.Hostname() == "" {
		return rc.trimDomain(urlStr)
	}
	return rc.trimDomain(u.Hostname())
}

// trimDomain trims domain.
func (rc *ResourceConfig) trimDomain(domain string) string {
	domainReturn := strings.TrimPrefix(
		strings.TrimPrefix(strings.TrimPrefix(domain, "http://"), "https://"), "www.",
	)
	return strings.Trim(domainReturn, "/ ")
}
