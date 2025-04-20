package rotations

import "github.com/nezbut/proxym"

// CompositeRotationLogicType is a type for composite rotation logic.
type CompositeRotationLogicType int

// CompositeRotationLogicType constants.
const (
	// RotationLogicAND returns true if all strategies return true.
	RotationLogicAND CompositeRotationLogicType = iota
	// RotationLogicOR returns true if any strategy returns true.
	RotationLogicOR
)

// CompositeRotation is a composite rotation strategy.
//
// It is used to determine if a proxy should be rotated based on multiple strategies.
type CompositeRotation struct {
	strategies []proxym.RotationStrategy
	logic      CompositeRotationLogicType
}

// NewCompositeRotationStrategy creates a new composite rotation strategy.
func NewCompositeRotationStrategy(
	logic CompositeRotationLogicType,
	strategies ...proxym.RotationStrategy,
) proxym.RotationStrategy {
	return &CompositeRotation{
		strategies: strategies,
		logic:      logic,
	}
}

// ShouldRotate returns true if the proxy should be rotated.
func (c *CompositeRotation) ShouldRotate(proxy *proxym.Proxy) bool {
	if len(c.strategies) == 0 {
		return false
	}

	for _, strategy := range c.strategies {
		result := strategy.ShouldRotate(proxy)

		if c.logic == RotationLogicOR && result {
			return true
		}

		if c.logic == RotationLogicAND && !result {
			return false
		}
	}
	return c.logic == RotationLogicAND
}
