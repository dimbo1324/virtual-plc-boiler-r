package domain

import "math"

type PIDController struct {
	Kp        float64
	Ki        float64
	Kd        float64
	prevError float64
	integral  float64
	minOut    float64
	maxOut    float64
}

func NewPID(kp, ki, kd float64) *PIDController {
	return &PIDController{
		Kp:     kp,
		Ki:     ki,
		Kd:     kd,
		minOut: 0.0,
		maxOut: 100.0,
	}
}
func (pid *PIDController) Update(setpoint, pv, dt float64) float64 {
	err := setpoint - pv
	P := pid.Kp * err
	pid.integral += err * dt
	if pid.Ki != 0 {
		if pid.Ki*pid.integral > pid.maxOut {
			pid.integral = pid.maxOut / pid.Ki
		} else if pid.Ki*pid.integral < pid.minOut {
			pid.integral = pid.minOut / pid.Ki
		}
	}
	I := pid.Ki * pid.integral
	derivative := (err - pid.prevError) / dt
	D := pid.Kd * derivative
	pid.prevError = err
	output := P + I + D
	return math.Max(pid.minOut, math.Min(pid.maxOut, output))
}
