package main

import (
	"testing"
)

func TestGetChildrenNumbers(t *testing.T) {
	testCases := []struct {
		min,
		max int
	}{
		{1, 10000},
		{10001, 20000},
		{20001, 30000},
		{40001, 50000},
		{70001, 80000},
		{90001, 100000},
	}
	p := 10
	for idx, tc := range testCases {
		childrenNumbers := getChildrenNumbers(int64(idx), tc.min, tc.max, p)
		for _, childrenNumber := range childrenNumbers {
			if tc.min > childrenNumber || tc.max < childrenNumber {
				t.Error("Expected to be within", tc.min, "and", tc.max, ", but got", childrenNumber, "for idx=", idx)
			}
		}
	}
}

func TestGetBatchStartAndEnd(t *testing.T) {
	testCases := []struct {
		level,
		batchNo,
		expStart,
		expEnd int
	}{
		{1, 1, 0, 999},
		{1, 4, 3000, 3999},
		{1, 10, 9000, 9999},
		{2, 1, 10000, 10999},
		{2, 10, 19000, 19999},
		{3, 1, 20000, 20999},
		{3, 4, 23000, 23999},
		{3, 10, 29000, 29999},
		{10, 1, 90000, 90999},
		{10, 10, 99000, 99999},
	}

	n = 10000
	for idx, tc := range testCases {
		actStart, actEnd := getBatchStartAndEnd(tc.level, tc.batchNo)
		if actStart != tc.expStart || actEnd != tc.expEnd {
			t.Errorf("Expected start=%d and end=%d, got start=%d and end=%d for idx=%d", tc.expStart, tc.expEnd, actStart, actEnd, idx)
		}
	}
}
