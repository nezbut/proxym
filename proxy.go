package proxym

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProxyPriority is a representation of a proxy priority in proxym.
type ProxyPriority uint

// Proxy priorities.
const (
	ProxyPriorityLow ProxyPriority = iota
	ProxyPriorityMedium
	ProxyPriorityHigh
)

// Proxy is a representation of a proxy in proxym.
//
// It has statistics and metadata can be useful for RotationStrategy and SelectStrategy.
//
// It can also be currently active or enabled/disabled.
type Proxy struct {
	url        *url.URL
	stats      *ProxyStats
	meta       *ProxyMetadata
	isActive   bool
	isDisabled bool
	mu         sync.RWMutex
}

// NewProxy creates a new Proxy.
func NewProxy(url *url.URL, meta *ProxyMetadata) *Proxy {
	if meta == nil {
		meta = &ProxyMetadata{}
	}
	return &Proxy{
		url:   url,
		meta:  meta,
		stats: &ProxyStats{},
	}
}

// NewProxyParsedStr creates a new Proxy from a string url.
//
// It returns an error if the url is invalid.
func NewProxyParsedStr(urlStr string, meta *ProxyMetadata) (*Proxy, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return NewProxy(u, meta), nil
}

// NewProxyStr creates a new Proxy from a string url.
//
// It panics if the url is invalid.
func NewProxyStr(urlStr string, meta *ProxyMetadata) *Proxy {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return NewProxy(u, meta)
}

// NewDirectConnection creates a proxy representing a direct connection.
func NewDirectConnection() *Proxy {
	return NewProxy(nil, nil)
}

// URL returns the proxy url.
func (p *Proxy) URL() *url.URL {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.url
}

// String returns the string representation of the proxy.
func (p *Proxy) String() string {
	u := p.URL()
	if u == nil {
		return "<not proxy url>"
	}
	return u.String()
}

// Disable marks the proxy as disabled.
func (p *Proxy) Disable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isDisabled = true
}

// Enable marks the proxy as enabled.
func (p *Proxy) Enable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isDisabled = false
}

// IsDisabled returns true if the proxy is disabled.
func (p *Proxy) IsDisabled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isDisabled
}

// activate marks the proxy as active.
func (p *Proxy) activate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isActive = true
}

// deactivate marks the proxy as inactive.
func (p *Proxy) deactivate() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isActive = false
}

// IsActive returns true if the proxy is active.
func (p *Proxy) IsActive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isActive
}

// IsDirect returns true if proxy represents a direct connection.
func (p *Proxy) IsDirect() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.url == nil
}

// Update is shorthand for Proxy.Stats().Update(response, err).
func (p *Proxy) Update(response *http.Response, err error) {
	p.Stats().Update(response, err)
}

// Stats returns the statistics of the proxy.
func (p *Proxy) Stats() *ProxyStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}

// Metadata returns the metadata of the proxy.
func (p *Proxy) Metadata() *ProxyMetadata {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.meta
}

// ProxyStats is a representation of a proxy statistics in proxym.
type ProxyStats struct {
	totalRequests uint
	successCount  uint
	errorCount    uint
	lastUsed      time.Time
	mu            sync.RWMutex
}

// TotalRequests returns the total requests of the proxy.
func (s *ProxyStats) TotalRequests() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalRequests
}

// SuccessCount returns the success count of the proxy.
func (s *ProxyStats) SuccessCount() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.successCount
}

// ErrorCount returns the error count of the proxy.
func (s *ProxyStats) ErrorCount() uint {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errorCount
}

// LastUsed returns the last used date of the proxy.
func (s *ProxyStats) LastUsed() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastUsed
}

// Update updates the proxy statistics at the expense of *http.Response and response error.
func (s *ProxyStats) Update(response *http.Response, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalRequests++

	if response != nil && err == nil {
		s.successCount++
	} else {
		s.errorCount++
	}

	s.lastUsed = time.Now()
}

// ProxyMetadata is a representation of a proxy metadata in proxym.
type ProxyMetadata struct {
	country   string
	priority  ProxyPriority
	expiresAt time.Time
	mu        sync.RWMutex
}

// NewProxyMetadata creates a new ProxyMetadata.
func NewProxyMetadata(country string, priority ProxyPriority, expiresAt time.Time) *ProxyMetadata {
	return &ProxyMetadata{
		country:   country,
		priority:  priority,
		expiresAt: expiresAt,
	}
}

// SetPriority sets the priority of the proxy.
func (m *ProxyMetadata) SetPriority(priority ProxyPriority) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.priority = priority
}

// Priority returns the priority of the proxy.
func (m *ProxyMetadata) Priority() ProxyPriority {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.priority
}

// SetCountry sets the country of the proxy.
func (m *ProxyMetadata) SetCountry(country string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.country = country
}

// Country returns the country of the proxy.
func (m *ProxyMetadata) Country() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.country
}

// SetExpiresAt sets the expiration date of the proxy.
func (m *ProxyMetadata) SetExpiresAt(expiresAt time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.expiresAt = expiresAt
}

// ExpiresAt returns the expiration date of the proxy.
func (m *ProxyMetadata) ExpiresAt() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.expiresAt
}
