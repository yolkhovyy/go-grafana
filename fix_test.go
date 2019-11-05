package main

import (
	"math"
	"testing"
)

// TestPassthrough validates that Fix returns the input data if no modifications are required.
func TestPassthrough(t *testing.T) {
	in := []Point{
		{0, 100},
		{8.4, 110},
		{0.999, 120},
	}
	got := Fix(in, 100, 130, 10)
	exp := []Point{
		{0, 100},
		{8.4, 110},
		{0.999, 120},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestApplyFromTo validates that Fix only returns points that lie within the from/to range provided.
func TestApplyFromTo(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{8.4, 110},
		{32.5, 120},
		{0.999, 130},
	}
	got := Fix(in, 110, 130, 10)
	exp := []Point{
		{8.4, 110},
		{32.5, 120},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestApplyFromToNaNs validates that Fix adds NaN's at the front or end as needed to match the requested range.
func TestApplyFromToNaNs(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{41.3, 110},
	}
	got := Fix(in, 90, 140, 10)
	exp := []Point{
		{math.NaN(), 90},
		{0.5, 100},
		{41.3, 110},
		{math.NaN(), 120},
		{math.NaN(), 130},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestApplyFromToNaNs validates that Fix honors from/to even when they are irregular (when they are not divisible by the interval).
func TestApplyIrregularFromToNaNs(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{41.3, 110},
	}
	got := Fix(in, 89, 131, 10)
	exp := []Point{
		{math.NaN(), 90},
		{0.5, 100},
		{41.3, 110},
		{math.NaN(), 120},
		{math.NaN(), 130},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestAdjustIrregularTimestamps validates that Fix adjusts timestamps to be divisible by the interval.
func TestAdjustIrregularTimestamps(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{8.4, 110},
		{32.5, 118},
		{0.999, 130},
		{41.3, 139},
		{41.9, 141},
	}
	got := Fix(in, 100, 160, 10)
	exp := []Point{
		{0.5, 100},
		{8.4, 110},
		{32.5, 120},
		{0.999, 130},
		{41.3, 140},
		{41.9, 150},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestInsertNaNs validates that Fix insert NaNs for missing points.
func TestInsertNaNs(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{32.5, 120},
		{0.999, 130},
		{41.3, 150},
	}
	got := Fix(in, 100, 160, 10)
	exp := []Point{
		{0.5, 100},
		{math.NaN(), 110},
		{32.5, 120},
		{0.999, 130},
		{math.NaN(), 140},
		{41.3, 150},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestFilterDuplicates validates that Fix filters out points that correspond to the same output timestamp.
func TestFilterDuplicates(t *testing.T) {
	in := []Point{
		{0.5, 100},
		{0.7, 100},
		{32.5, 110},
		{0.9, 120},
		{0.8, 125},
		{41.3, 130},
		{41.3, 140},
		{41.2, 141},
		{41.3, 142},
		{41.4, 149},
	}
	got := Fix(in, 100, 160, 10)
	exp := []Point{
		{0.5, 100},
		{32.5, 110},
		{0.9, 120},
		{0.8, 130},
		{41.3, 140},
		{41.2, 150},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestIrregularFromToIrregularTimestamps validates that Fix correctly handles the combination
// of irregular from/to parameters with irregular input timestamps (irregular meaning not divisible by interval)
// note that the from/to filtering should be applied to the adjusted timestamps.
func TestIrregularFromToIrregularTimestamps(t *testing.T) {
	in := []Point{
		{0.5, 105},
		{1, 120},
		{1.5, 130},
		{2, 135},
		{3, 139},
	}
	got := Fix(in, 108, 138, 10)
	exp := []Point{
		{0.5, 110},
		{1, 120},
		{1.5, 130},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

// TestAllAtOnce validates that Fix can handle all the above conditions all at once.
func TestAllAtOnce(t *testing.T) {
	in := []Point{
		{0.5, 105},
		{41.3, 111},
		{0, 130},
		{1, 150},
		{2, 151},
		{3, 152},
		{4, 159},
		{4, 160},
		{5, 160},
		{4, 174},
		{5, 175},
		{6, 185},
	}
	got := Fix(in, 81, 186, 10)
	exp := []Point{
		{math.NaN(), 90},
		{math.NaN(), 100},
		{0.5, 110},
		{41.3, 120},
		{0, 130},
		{math.NaN(), 140},
		{1, 150},
		{2, 160},
		{math.NaN(), 170},
		{4, 180},
	}

	if !equal(exp, got) {
		t.Fatalf("output mismatch:\nexpected:\n%v\ngot:\n%v", exp, got)
	}
}

func equal(exp, got []Point) bool {
	if len(exp) != len(got) {
		return false
	}
	for i, pgot := range got {
		pexp := exp[i]
		if math.IsNaN(pgot.Val) != math.IsNaN(pexp.Val) {
			return false
		}
		if !math.IsNaN(pgot.Val) && pgot.Val != pexp.Val {
			return false
		}
		if pgot.Ts != pexp.Ts {
			return false
		}
	}
	return true
}
