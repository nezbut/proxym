package proxym

// RotationStrategy is an interface for proxy rotation strategies.
// It is used to determine if a proxy should be rotated.
type RotationStrategy interface {
	// ShouldRotate returns true if the proxy should be rotated.
	ShouldRotate(proxy *Proxy) bool
}
