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
	}
	if limits != nil {
		t.Errorf("Non-nil limits returned, expected nil for uncached account.")
	}

}
