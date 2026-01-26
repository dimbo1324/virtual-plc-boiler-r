package physics

import (
	"math"
	"testing"
	"time"
)

// floatEquals checks if two floats are equal within a small epsilon.
// Floating point math is rarely exact, so we compare ranges.
func floatEquals(a, b float64) bool {
	const epsilon = 1e-9
	return math.Abs(a-b) < epsilon
}

func TestSolarIrradiance(t *testing.T) {
	// 1. Define the test cases (Table-Driven approach)
	tests := []struct {
		name           string
		checkTime      time.Time // Input time
		peak           float64
		sunrise        float64
		sunset         float64
		expectedResult float64 // Approximate expectation
		shouldBeZero   bool    // If true, strictly check for 0.0
	}{
		{
			name:         "Night time (Midnight)",
			checkTime:    time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			peak:         1000,
			sunrise:      6.0,
			sunset:       20.0,
			shouldBeZero: true,
		},
		{
			name: "Solar Noon (Peak)",
			// 13:00 is exactly between 06:00 and 20:00 (14 hour day, 7 hours in)
			checkTime:      time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC),
			peak:           1000,
			sunrise:        6.0,
			sunset:         20.0,
			expectedResult: 1000.0, // Sin(Pi/2) = 1
		},
		{
			name:         "Sunrise boundary",
			checkTime:    time.Date(2023, 1, 1, 6, 0, 0, 0, time.UTC),
			peak:         1000,
			sunrise:      6.0,
			sunset:       20.0,
			shouldBeZero: true,
		},
	}

	// 2. Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SolarIrradiance(tc.checkTime, tc.peak, tc.sunrise, tc.sunset)

			if tc.shouldBeZero {
				if got != 0.0 {
					t.Errorf("SolarIrradiance() = %v, want 0.0", got)
				}
			} else {
				// Allow small margin of error for non-zero calculations
				if math.Abs(got-tc.expectedResult) > 0.001 {
					t.Errorf("SolarIrradiance() = %v, want %v", got, tc.expectedResult)
				}
			}
		})
	}
}

func TestSolarPanelPower(t *testing.T) {
	type args struct {
		irradiance float64
		area       float64
		efficiency float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Standard calculation",
			args: args{irradiance: 1000, area: 2, efficiency: 0.2},
			want: 400.0, // 1000 * 2 * 0.2 = 400
		},
		{
			name: "Zero irradiance",
			args: args{irradiance: 0, area: 5, efficiency: 0.2},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SolarPanelPower(tt.args.irradiance, tt.args.area, tt.args.efficiency); !floatEquals(got, tt.want) {
				t.Errorf("SolarPanelPower() = %v, want %v", got, tt.want)
			}
		})
	}
}
