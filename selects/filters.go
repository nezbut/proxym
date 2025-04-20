package selects

import "github.com/nezbut/proxym"

// RemoveActiveProxyFilter filters and removes the active proxy.
type RemoveActiveProxyFilter struct{}

// Filter returns the filtered list of proxies.
func (f RemoveActiveProxyFilter) Filter(proxies []*proxym.Proxy) []*proxym.Proxy {
	result := make([]*proxym.Proxy, 0, len(proxies))
	for _, p := range proxies {
		if !p.IsActive() {
			result = append(result, p)
		}
	}
	return result
}

// RemoveDisabledFilter filters and removes the disabled proxies.
type RemoveDisabledFilter struct{}

// Filter returns the filtered list of proxies.
func (f RemoveDisabledFilter) Filter(proxies []*proxym.Proxy) []*proxym.Proxy {
	result := make([]*proxym.Proxy, 0, len(proxies))
	for _, p := range proxies {
		if !p.IsDisabled() {
			result = append(result, p)
		}
	}
	return result
}
