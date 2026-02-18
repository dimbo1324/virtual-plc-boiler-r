package domain

// PIDController хранит коэффициенты и состояние ошибки
type PIDController struct {
	Kp, Ki, Kd float64
}

// NewPID создает новый регулятор
func NewPID(kp, ki, kd float64) *PIDController {
	return &PIDController{Kp: kp, Ki: ki, Kd: kd}
}

// Update считает управляющее воздействие (0-100%)
func (pid *PIDController) Update(setpoint, pv, dt float64) float64 {
	return 0
}
