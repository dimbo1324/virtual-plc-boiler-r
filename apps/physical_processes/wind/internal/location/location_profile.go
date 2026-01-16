package location

import "time"

// Состояние Среды
type CurrentWeather struct {
	t    float64       // Температура (для расчета плотности).
	p    float64       // Давление (для расчета плотности).
	time time.Duration // Временная метка (для определения времени суток).
}
