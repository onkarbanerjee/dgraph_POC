package main

import (
	"context"
	"testing"
)

func TestServiceToAlarmCorrelation(t *testing.T) {
	// t.Parallel()

	dgClient, cancel := getDgraphClient("100.64.1.15:9080")
	defer cancel()

	testCases := []struct {
		m   int
		fre string
	}{
		// {2, "FRE11234"},
		{3, "FRE21234"},
		// {4, "FRE31234"},
		// {5, "FRE41234"},
		// {6, "FRE51234"},
		// {7, "FRE61234"},
		// {8, "FRE71234"},
		// {9, "FRE81234"},
		// {10, "FRE91234"},
	}

	for _, tc := range testCases {
		res, err := serviceToAlarmCorrelation(context.Background(), dgClient, tc.fre)
		if err != nil {
			t.Error("Expected no err, got", err, " at m=", tc.m)
		}
		t.Log("\nAt m=", tc.m, "time taken =", res.GetLatency().String())
	}
}

// func timeTrack(start time.Time, name string) string {
// 	elapsed := time.Since(start)
// 	return fmt.Sprintf("%s took %s", name, elapsed)
// }
