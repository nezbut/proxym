package proxym

import (
	"net/http"
)

// ProxyTransport is http.RoundTripper that first receives the response through the base transport
// Then updates the proxy data.
//
// The base transport must receive a proxy via ProxySelector for requests.
type ProxyTransport struct {
	pm            ProxyManager
	baseTransport http.RoundTripper
}

// NewProxyTransport returns a new ProxyTransport.
func NewProxyTransport(pm ProxyManager, baseTransport http.RoundTripper) *ProxyTransport {
	return &ProxyTransport{pm: pm, baseTransport: baseTransport}
}

// RoundTrip calls the base transport and updates the proxy data.
func (pt *ProxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := pt.baseTransport.RoundTrip(req)
	proxy := pt.pm.LastUsed()
	if proxy != nil {
		proxy.Update(resp, err)
	}
	return resp, err
}

// NewClient returns a new http.Client with a ProxyTransport and with a cloned http.DefaultTransport.
func NewClient(pm ProxyManager) *http.Client {
	cloned, _ := CloneRoundTripperWithProxySelector(pm, http.DefaultTransport)
	return &http.Client{
		Transport: NewProxyTransport(pm, cloned),
	}
}

// PatchClient patches the http.Client with a ProxyTransport and with a cloned client.Transport.
//
// Call this function in the application initialization, as this function is not thread-safe.
func PatchClient(client *http.Client, pm ProxyManager) error {
	if client.Transport == nil {
		cloned, _ := CloneRoundTripperWithProxySelector(pm, http.DefaultTransport)
		client.Transport = NewProxyTransport(pm, cloned)
	} else {
		cloned, err := CloneRoundTripperWithProxySelector(pm, client.Transport)
		if err != nil {
			return err
		}
		client.Transport = NewProxyTransport(pm, cloned)
	}
	return nil
}
