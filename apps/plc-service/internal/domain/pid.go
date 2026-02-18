package domain

import "math"

// PIDController хранит коэффициенты и состояние ошибки
type PIDController struct {
	Kp        float64 // Пропорциональный коэффициент
	Ki        float64 // Интегральный коэффициент
	Kd        float64 // Дифференциальный коэффициент
	prevError float64 // Ошибка на предыдущем шаге
	integral  float64 // Накопленная сумма ошибок
	minOut    float64 // Минимальный выход (0%)
	maxOut    float64 // Максимальный выход (100%)
}

// NewPID создает новый регулятор
func NewPID(kp, ki, kd float64) *PIDController {
	return &PIDController{
		Kp:     kp,
		Ki:     ki,
		Kd:     kd,
		minOut: 0.0,
		maxOut: 100.0,
	}
}

/*
Update рассчитывает управляющее воздействие.
setpoint - целевое значение (например, 60 бар)
pv - текущее значение (process variable, текущее давление)
dt - время в секундах, прошедшее с прошлого расчета
*/
func (pid *PIDController) Update(setpoint, pv, dt float64) float64 {
	err := setpoint - pv

	P := pid.Kp * err

	pid.integral += err * dt
	I := pid.Ki * pid.integral

	derivative := (err - pid.prevError) / dt
	D := pid.Kd * derivative

	pid.prevError = err
	output := P + I + D

	return math.Max(pid.minOut, math.Min(pid.maxOut, output))
}
