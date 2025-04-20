package rotations

import "github.com/nezbut/proxym"

// RequestLimitedRotation is a rotation strategy that returns true
// if the total number of requests is greater than or equal to a limit.
type RequestLimitedRotation struct {
	limit uint
}

// NewRequestLimitedRotation returns a new RequestLimitedRotation.
func NewRequestLimitedRotation(limit uint) proxym.RotationStrategy {
	return &RequestLimitedRotation{limit: limit}
}

// ShouldRotate returns true if the proxy need is rotated.
func (r *RequestLimitedRotation) ShouldRotate(proxy *proxym.Proxy) bool {
	return proxy.Stats().TotalRequests() >= r.limit
}
