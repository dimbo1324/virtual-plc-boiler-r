package physics

import (
	"math"
	"testing"
)

func TestWeibullWindSpeed(t *testing.T) {
	tests := []struct {
		name  string
		u     float64 // Random input [0, 1)
		shape float64
		scale float64
		want  float64
	}{
		{
			name:  "Zero input (Calm)",
			u:     0.0,
			shape: 2.0,
			scale: 10.0,
			want:  0.0, // ln(1-0) = ln(1) = 0
		},
		{
			name: "Standard input (0.632)",
			// If u = 1 - e^(-1) approx 0.632, then -ln(1-u) = 1. Result should be scale.
			u:     1.0 - math.Exp(-1.0),
			shape: 2.0,
			scale: 10.0,
			want:  10.0,
		},
		{
			name:  "Out of bounds input (Negative)",
			u:     -0.5,
			shape: 2.0,
			scale: 10.0,
			want:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeibullWindSpeed(tt.u, tt.shape, tt.scale)
			// Check with tolerance
			if math.Abs(got-tt.want) > 1e-5 {
				t.Errorf("WeibullWindSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWindTurbinePower(t *testing.T) {
	type args struct {
		rho       float64
		area      float64
		windSpeed float64
		cp        float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Standard calculation",
			args: args{
				rho:       1.225,
				area:      10.0,
				windSpeed: 10.0, // v^3 = 1000
				cp:        0.4,
			},
			// 0.5 * 1.225 * 10 * 1000 * 0.4 = 2450
			want: 2450.0,
		},
		{
			name: "No wind",
			args: args{rho: 1.225, area: 10, windSpeed: 0, cp: 0.4},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WindTurbinePower(tt.args.rho, tt.args.area, tt.args.windSpeed, tt.args.cp)
			if math.Abs(got-tt.want) > 1e-5 {
				t.Errorf("WindTurbinePower() = %v, want %v", got, tt.want)
			}
		})
	}
}
