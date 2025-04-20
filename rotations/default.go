package rotations

import "github.com/nezbut/proxym"

// DefaultRotationStrategy returns the default rotation strategy.
//
// It returns true if the proxy is disabled or if the error proxy is greater than or equal to 1.
func DefaultRotationStrategy() proxym.RotationStrategy {
	return NewCompositeRotationStrategy(
		RotationLogicOR,
		OnlyEnabledRotation{},
		NewErrorThresholdRotation(1),
	)
}
