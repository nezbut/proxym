package selects

import "github.com/nezbut/proxym"

// SelectFilter is an interface for proxy selection strategies filters.
//
// It is used to filter the list of proxies before selecting a proxy.
type SelectFilter interface {
	// Filter returns the filtered list of proxies.
	Filter(proxies []*proxym.Proxy) []*proxym.Proxy
}

// FilteredSelectProvider is a provider that first gets the proxies from the source provider
// filters them and then returns them.
type FilteredSelectProvider struct {
	sourceProvider proxym.SelectStrategyProxyProvider
	filters        []SelectFilter
}

// NewFilteredSelectProvider creates a new FilteredSelectProvider.
func NewFilteredSelectProvider(
	sourceProvider proxym.SelectStrategyProxyProvider,
	filters ...SelectFilter,
) proxym.SelectStrategyProxyProvider {
	return &FilteredSelectProvider{
		sourceProvider: sourceProvider,
		filters:        filters,
	}
}

// NewFilteredSelectFactory creates a new proxym.SelectStrategyFactory
// that injects selects.FilteredSelectProvider into proxym.SelectStrategy along with some source provider and filters.
func NewFilteredSelectFactory(
	selectFactory proxym.SelectStrategyFactory, // The select strategy factory to use.
	filters ...SelectFilter, // The filters
) proxym.SelectStrategyFactory {
	return func(sourceProvider proxym.SelectStrategyProxyProvider) proxym.SelectStrategy {
		return selectFactory(NewFilteredSelectProvider(sourceProvider, filters...))
	}
}

// GetProxies returns the filtered list of proxies.
func (f *FilteredSelectProvider) GetProxies() []*proxym.Proxy {
	proxies := f.sourceProvider.GetProxies()

	for _, filter := range f.filters {
		proxies = filter.Filter(proxies)
		if len(proxies) == 0 {
			return proxies
		}
	}
	return proxies
}
