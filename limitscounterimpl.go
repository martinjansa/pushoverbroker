package main

import (
	"errors"
	"sync"
)

type limitsCache map[string]*Limits

// LimitsCounterImpl implements the LImitsCounter interface
type LimitsCounterImpl struct {
	limitsCache      limitsCache
	limitsCacheMutex sync.Mutex
}

// SetLimits stores the current limits values for the give account. Should be called after a successful connection to Pushover servers
func (l *LimitsCounterImpl) SetLimits(accountToken string, limits *Limits) error {

	// lock the mutex
	l.limitsCacheMutex.Lock()
	defer l.limitsCacheMutex.Unlock()

	// store the limits value for the account
	l.limitsCache[accountToken] = limits
	return nil
}

// DecrementLimits decrements the limit of the available messages for the given account
func (l *LimitsCounterImpl) DecrementLimits(accountToken string) error {

	// lock the mutex
	l.limitsCacheMutex.Lock()
	defer l.limitsCacheMutex.Unlock()

	// seach the account in the map
	limits, exists := l.limitsCache[accountToken]
	if !exists {
		return errors.New("the requested account does not exist in the cache")
	}

	// if there are still some remaining messages in the limits
	if limits.remaining > 0 {

		// decrement the limits
		limits.remaining--
	}

	return nil
}

// GetLimits returns the current limits or nil, if not known yet
func (l *LimitsCounterImpl) GetLimits(accountToken string) (*Limits, error) {

	// lock the mutex
	l.limitsCacheMutex.Lock()
	defer l.limitsCacheMutex.Unlock()

	// seach the account in the map
	limits, exists := l.limitsCache[accountToken]
	if !exists {

		// return empty limits and no error
		return nil, nil
	}

	// return limits and no error
	return limits, nil
}

// NewLimitsCounterImpl creates a new limits counter instance
func NewLimitsCounterImpl() *LimitsCounterImpl {
	lc := new(LimitsCounterImpl)
	lc.limitsCache = make(limitsCache)
	return lc
}
