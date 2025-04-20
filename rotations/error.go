package rotations

import "github.com/nezbut/proxym"

// ErrorThresholdRotation is a rotation strategy that returns true
// if the error proxy is greater than or equal to a threshold.
type ErrorThresholdRotation struct {
	threshold uint
}

// NewErrorThresholdRotation returns a new ErrorThresholdRotation.
func NewErrorThresholdRotation(threshold uint) proxym.RotationStrategy {
	return &ErrorThresholdRotation{threshold: threshold}
}

// ShouldRotate returns true if the proxy need is rotated.
func (e *ErrorThresholdRotation) ShouldRotate(proxy *proxym.Proxy) bool {
	return proxy.Stats().ErrorCount() >= e.threshold
}
