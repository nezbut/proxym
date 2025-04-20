package rotations

import "github.com/nezbut/proxym"

// OnlyEnabledRotation is a rotation strategy that returns true if the proxy is disabled.
//
// Returns true, which means that should rotate proxy if proxy is disabled.
type OnlyEnabledRotation struct{}

// ShouldRotate returns true if the proxy need is rotated.
func (o OnlyEnabledRotation) ShouldRotate(proxy *proxym.Proxy) bool {
	return proxy.IsDisabled()
}
