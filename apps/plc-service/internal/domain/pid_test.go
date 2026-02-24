package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPIDController_Update_Proportional(t *testing.T) {
	pid := NewPID(2.0, 0.0, 0.0)

	tests := []struct {
		name     string
		setpoint float64
		pv       float64
		dt       float64
		expected float64
	}{
		{
			name:     "Zero error yields zero output",
			setpoint: 60.0,
			pv:       60.0,
			dt:       1.0,
			expected: 0.0,
		},
		{
			name:     "Positive error yields positive output",
			setpoint: 60.0,
			pv:       50.0,
			dt:       1.0,
			expected: 20.0,
		},
		{
			name:     "Clamping Max (Output > 100 should be 100)",
			setpoint: 100.0,
			pv:       0.0,
			dt:       1.0,
			expected: 100.0,
		},
		{
			name:     "Clamping Min (Output < 0 should be 0)",
			setpoint: 50.0,
			pv:       60.0,
			dt:       1.0,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := pid.Update(tt.setpoint, tt.pv, tt.dt)
			assert.InDelta(t, tt.expected, output, 0.001, "Output mismatch")
		})
	}
}

func TestPIDController_Update_Integral(t *testing.T) {
	pid := NewPID(0.0, 1.0, 0.0)

	out1 := pid.Update(60.0, 50.0, 1.0)
	assert.InDelta(t, 10.0, out1, 0.001)

	out2 := pid.Update(60.0, 50.0, 1.0)
	assert.InDelta(t, 20.0, out2, 0.001)
}
