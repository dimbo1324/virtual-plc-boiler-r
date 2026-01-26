package physics

import (
	"math"
	"time"
)

// SolarIrradiance calculates the current solar irradiance in Watts per square meter (W/m²)
// based on the time of day, simulating the sun's path with a sinusoidal model.
//
// This function assumes a simplified day-night cycle and uses a sine wave to model
// the irradiance curve, peaking at noon. It returns 0 during nighttime hours or if
// the calculated value is negative (due to floating-point precision).
//
// Parameters:
//   - t: The current time (time.Time object).
//   - peakIrradiance: The maximum irradiance at noon (W/m²), which can vary based on
//     weather conditions, cloud cover, or location-specific factors.
//   - sunriseHour: The hour of sunrise as a float (e.g., 6.0 for 6:00 AM).
//   - sunsetHour: The hour of sunset as a float (e.g., 20.0 for 8:00 PM).
//
// Returns:
//
//	The calculated solar irradiance (W/m²). Returns 0 if outside daylight hours or
//	if the computation yields a negative value.
func SolarIrradiance(t time.Time, peakIrradiance, sunriseHour, sunsetHour float64) float64 {
	// Calculate the current hour as a fractional value (e.g., 14.5 for 2:30 PM)
	currentHour := float64(t.Hour()) + float64(t.Minute())/60.0

	// If it's before sunrise or after sunset, return 0 irradiance (nighttime)
	if currentHour < sunriseHour || currentHour > sunsetHour {
		return 0.0
	}

	// Compute the duration of the daylight period in hours
	dayDuration := sunsetHour - sunriseHour

	// Normalize the time from sunrise to a radian value between 0 and π for the sine function
	// This maps the daylight period to a half-sine wave
	radians := (currentHour - sunriseHour) / dayDuration * math.Pi

	// Simulate the sun's parabolic irradiance curve using the sine function
	irradiance := peakIrradiance * math.Sin(radians)

	// Ensure no negative values are returned (safeguard for precision issues)
	if irradiance < 0 {
		return 0
	}

	return irradiance
}

// SolarPanelPower computes the power output of a solar panel in Watts (W) based on
// the incident irradiance, panel area, and efficiency.
//
// This is a straightforward calculation using the formula: Power = Irradiance × Area × Efficiency.
// It assumes uniform irradiance across the panel and does not account for factors like
// temperature, angle of incidence, or degradation over time.
//
// Parameters:
//   - irradiance: The incident solar irradiance (W/m²).
//   - area: The surface area of the solar panel (m²).
//   - efficiency: The efficiency of the panel as a decimal (e.g., 0.15 for 15%, typically
//     ranging from 0.15 to 0.22 for standard panels).
//
// Returns:
//
//	The power output in Watts (W).
func SolarPanelPower(irradiance, area, efficiency float64) float64 {
	return irradiance * area * efficiency
}
