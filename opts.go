package proxym

// ProxyManagerImplOption is option for ProxyManagerImpl.
type ProxyManagerImplOption func(*ProxyManagerImpl)

// WithProxies sets proxies to the ProxyManagerImpl.
func WithProxies(proxies ...*Proxy) ProxyManagerImplOption {
	return func(pm *ProxyManagerImpl) {
		pm.proxies = proxies
	}
}

// WithResources sets resources to the ProxyManagerImpl.
func WithResources(resources ...*ResourceConfig) ProxyManagerImplOption {
	return func(pm *ProxyManagerImpl) {
		pm.resources = resources
	}
}

// WithRotationStrategy sets rotation strategy to the ProxyManagerImpl.
func WithRotationStrategy(strategy RotationStrategy) ProxyManagerImplOption {
	return func(pm *ProxyManagerImpl) {
		pm.rotationStrategy = strategy
	}
}

// WithSelectStrategy sets select strategy from factory to the ProxyManagerImpl.
func WithSelectStrategy(factory SelectStrategyFactory) ProxyManagerImplOption {
	return func(pm *ProxyManagerImpl) {
		pm.selectStrategy = factory(pm)
	}
}

// ResourceConfigOption is option for ResourceConfig.
type ResourceConfigOption func(*ResourceConfig)

// WithResourceProxies sets proxies to the ResourceConfig.
func WithResourceProxies(proxies ...*Proxy) ResourceConfigOption {
	return func(rc *ResourceConfig) {
		rc.proxies = proxies
	}
}

// WithResourceSelectStrategy sets select strategy from factory to the ResourceConfig.
func WithResourceSelectStrategy(factory SelectStrategyFactory) ResourceConfigOption {
	return func(rc *ResourceConfig) {
		rc.selectStrategy = factory(rc)
	}
}

// WithResourceRotationStrategy sets rotation strategy to the ResourceConfig.
func WithResourceRotationStrategy(strategy RotationStrategy) ResourceConfigOption {
	return func(rc *ResourceConfig) {
		rc.rotationStrategy = strategy
	}
}

// WithDomain sets domain to the ResourceConfig.
func WithDomain(domain string) ResourceConfigOption {
	return func(rc *ResourceConfig) {
		rc.domain = domain
	}
}

// WithIgnoreSubdomains sets ignore subdomains to the ResourceConfig.
//
// If ignore is true, then it will ignore subdomains in the comparison of the domain.
func WithIgnoreSubdomains(ignore bool) ResourceConfigOption {
	return func(rc *ResourceConfig) {
		rc.notIgnoreSubdomains = !ignore
	}
}
