package proxym

import (
	"net/http"
	"net/url"
)

// ProxySelector is a function that returns the next available proxy url by request.
type ProxySelector func(*http.Request) (*url.URL, error)

// ProxySelectorSetter is an interface that allows to set a ProxySelector to a http.RoundTripper.
type ProxySelectorSetter interface {
	// WithProxySelector sets the ProxySelector to the http.RoundTripper.
	WithProxySelector(selector ProxySelector) http.RoundTripper
}

// CloneRoundTripperWithProxySelector returns a cloned http.RoundTripper with a ProxySelector.
//
// If the http.RoundTripper implementation is not supported, it returns an ErrUnsupportedRoundTripperImpl.
// Supported http.RoundTripper: http.Transport and ProxySelectorSetter.
func CloneRoundTripperWithProxySelector(pm ProxyManager, rt http.RoundTripper) (http.RoundTripper, error) {
	switch t := rt.(type) {
	case *http.Transport:
		cloned := t.Clone()
		cloned.Proxy = GetProxySelector(pm)
		return cloned, nil
	case ProxySelectorSetter:
		return t.WithProxySelector(GetProxySelector(pm)), nil
	default:
		return nil, ErrUnsupportedRoundTripperImpl
	}
}

// GetProxySelector returns a ProxySelector that uses the ProxyManager to get the next available proxy.
func GetProxySelector(pm ProxyManager) ProxySelector {
	return func(req *http.Request) (*url.URL, error) {
		proxy, err := pm.GetNextProxy(req.URL.Hostname())
		if err != nil {
			return nil, err
		}
		if proxy.IsDisabled() {
			return nil, ErrProxyNotAvailable
		}
		return proxy.url, nil
	}
}
