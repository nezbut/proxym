package selects

import "github.com/nezbut/proxym"

// DefaultSelectStrategy returns the default select strategy.
//
// It returns a RandomSelect with RemoveActiveProxyFilter and RemoveDisabledFilter.
func DefaultSelectStrategy() proxym.SelectStrategyFactory {
	return NewFilteredSelectFactory(
		NewRandomSelect,
		RemoveActiveProxyFilter{},
		RemoveDisabledFilter{},
	)
}
