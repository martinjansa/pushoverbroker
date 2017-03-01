package main

import "errors"

// LimitsCounterImpl implements the LImitsCounter interface
type LimitsCounterImpl struct {
	// limitsCache map[string]*Limits
	// limitsCacheMutex sync.Mutex
}

// SetLimits stores the current limits values for the give account. Should be called after a successful connection to Pushover servers
func (l *LimitsCounterImpl) SetLimits(accountToken string, limits *Limits) error {
	return errors.New("not implemented")
}

// DecrementLimits decrements the limit of the available messages for the given account
func (l *LimitsCounterImpl) DecrementLimits(accountToken string) error {
	return errors.New("not implemented")
}

// GetLimits returns the current limits or nil, if not known yet
func (l *LimitsCounterImpl) GetLimits(accountToken string) (*Limits, error) {
	return nil, errors.New("not implemented")
}

func NewLimitsCounterImpl() *LimitsCounterImpl {
	lc := new(LimitsCounterImpl)
	return lc
}
