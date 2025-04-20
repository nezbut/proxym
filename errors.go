package proxym

import "errors"

// Errors.
var (
	ErrProxyNotAvailable           = errors.New("proxy not available")
	ErrUnsupportedRoundTripperImpl = errors.New("unsupported round tripper implementation")
	ErrResourceNotFound            = errors.New("resource not found")
	ErrEmptyProxyList              = errors.New("empty proxy list in proxy manager")
	ErrFailedSelectProxy           = errors.New("failed select proxy in select strategy")
)
