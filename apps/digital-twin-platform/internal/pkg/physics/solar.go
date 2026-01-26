package physics

import (
	"math"
	"time"
)

func SolarIrradiance(t time.Time, peakIrradiance, sunriseHour, sunsetHour float64) float64 {
	currentHour := float64(t.Hour()) + float64(t.Minute())/60.0
	if currentHour < sunriseHour || currentHour > sunsetHour {
		return 0.0
	}
	dayDuration := sunsetHour - sunriseHour
	radians := (currentHour - sunriseHour) / dayDuration * math.Pi
	irradiance := peakIrradiance * math.Sin(radians)
	if irradiance < 0 {
		return 0
	}
	return irradiance
}

func SolarPanelPower(irradiance, area, efficiency float64) float64 {
	return irradiance * area * efficiency
}
