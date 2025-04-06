package handlers

import "testing"

func TestCalculateMean (t *testing.T) {
	testCases := []struct {
		name string
		input []float64
		want float64
	}{
		{
			name: "Empty Slice",
			input: []float64{},
			want: 0.0,
		},
		{
			name: "Single Element",
			input: []float64 {5.5},
			want: 5.50,
		},
		{
			name: "Multiple Positive Integers",
			input: []float64 {1.0, 2.0, 3.0},
			want: 2.00,
		},
		{
			name: "Mixed Positive and Negative",
			input: []float64 {-1.0, 1.0, 3.0, 5.0},
			want: 2.00,
		},
		{
			name: "Needs Rounding Up",
			input: []float64 {1.111, 2.222, 3.333},
			want: 2.22,
		},
		{
			name: "Needs Rounding Up (midpoint)",
			input: []float64 {1.115, 2.225},
			want: 1.67,
		},
		{
			name: "Needs Rounding Down",
			input: []float64 {10.0, 20.0, 35.0},
			want: 21.67,
		},
		{
			name: "Already Two Decimal",
			input: []float64 {1.23, 4.56, 7.89},
			want: 4.56,
		},
		{
			name: "Zeros",
			input: []float64 {0.0, 0.0, 0.0},
			want: 0.00,
		},
	}

	for _, tc := range testCases {
		got := calculateMean(tc.input)
		if got != tc.want {
			t.Errorf("calculateMean(%v) == %v, want %v", tc.input, got, tc.want)
		}
	}
}