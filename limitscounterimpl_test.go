package main

import "testing"

func TestLimitsCounterShouldGiveNilLimitsOnUncachedAccount(t *testing.T) {

	// GIVEN
	limitsCounterImpl := NewLimitsCounterImpl()

	// WHEN
	limits, err := limitsCounterImpl.GetLimits("uncacheckaccount")

	// THEN
	if err != nil {
		t.Errorf("limits getting failed with error %s, expected no error", err)
		return
	}
	if limits != nil {
		t.Errorf("Non-nil limits returned, expected nil for uncached account.")
		return
	}

}

func TestLimitsCounterShouldGiveOriginalValuesOnCachedAccount(t *testing.T) {

	// GIVEN
	limitsCounterImpl := NewLimitsCounterImpl()
	limitsCounterImpl.SetLimits("accountA", &Limits{limit: 1000, remaining: 500, reset: 123456789})

	// WHEN
	limits, err := limitsCounterImpl.GetLimits("accountA")

	// THEN
	if err != nil {
		t.Errorf("limits getting failed with error %s, expected no error", err)
		return
	}
	if limits == nil {
		t.Errorf("No limits returned, expected value.")
		return
	}
	if limits.limit != 1000 || limits.remaining != 500 || limits.reset != 123456789 {
		t.Errorf("Limits {%d, %d, %d} returned, expected {%d, %d, %d}.", limits.limit, limits.remaining, limits.reset, 1000, 500, 123456789)
		return
	}

}

func TestLimitsCounterShouldGiveDecrementedValuesOnCachedAccount(t *testing.T) {

	// GIVEN
	limitsCounterImpl := NewLimitsCounterImpl()
	limitsCounterImpl.SetLimits("accountA", &Limits{limit: 1000, remaining: 500, reset: 123456789})

	// WHEN
	decErr := limitsCounterImpl.DecrementLimits("accountA")
	limits, err := limitsCounterImpl.GetLimits("accountA")

	// THEN
	if decErr != nil {
		t.Errorf("limits decrementing failed with error %s, expected no error", decErr)
		return
	}
	if err != nil {
		t.Errorf("limits getting failed with error %s, expected no error", err)
		return
	}
	if limits == nil {
		t.Errorf("No limits returned, expected value.")
		return
	}
	if limits.limit != 1000 || limits.remaining != 500 || limits.reset != 123456789 {
		t.Errorf("Limits {%d, %d, %d} returned, expected {%d, %d, %d}.", limits.limit, limits.remaining, limits.reset, 1000, 500, 123456789)
		return
	}

}
