package main

// LimitsCounter represents an interface for caching the limits for account
type LimitsCounter interface {

	// SetLimits stores the current limits values for the give account. Should be called after a successful connection to Pushover servers
	SetLimits(accountToken string, limits *Limits) error

	// DecrementLimits decrements the limit of the available messages for the given account
	DecrementLimits(accountToken string) error

	// GetLimits returns the current limits or nil, if not known yet
	GetLimits(accountToken string) (*Limits, error)
}
