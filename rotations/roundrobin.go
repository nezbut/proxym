package rotations

import "github.com/nezbut/proxym"

// RoundRobinRotation is a rotation strategy that always returns true so that the proxies change on each request.
type RoundRobinRotation struct{}

// ShouldRotate returns true if the proxy need is rotated.
func (s RoundRobinRotation) ShouldRotate(_ *proxym.Proxy) bool {
	return true
}
