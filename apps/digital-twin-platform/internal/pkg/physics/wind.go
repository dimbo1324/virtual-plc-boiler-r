package physics

import "math"

// WeibullWindSpeed generates a wind speed value based on the Weibull distribution,
// commonly used to model wind speed variations in renewable energy simulations.
//
// This function employs the inverse cumulative distribution function (CDF) of the
// Weibull distribution to produce a realistic wind speed from a uniform random variable.
// It's particularly useful for stochastic modeling in wind energy systems, where wind
// speeds are not constant but follow a probabilistic pattern.
//
// Parameters:
//   - u: A uniform random number between 0.0 (exclusive) and 1.0 (exclusive), typically
//     generated from a random source in the environment.
//   - shape: The shape parameter (k) of the Weibull distribution (commonly around 2.0
//     for wind speeds, influencing the distribution's skewness).
//   - scale: The scale parameter (lambda) of the Weibull distribution (approximates the
//     mean wind speed in m/s, adjusting the distribution's spread).
//
// Returns:
//
//	The generated wind speed in meters per second (m/s). Returns 0.0 if u is invalid
//	(≤ 0 or ≥ 1).
func WeibullWindSpeed(u, shape, scale float64) float64 {
	if u <= 0 || u >= 1 {
		return 0.0
	}
	// Apply the inverse Weibull CDF formula: scale * (-ln(1 - u))^(1 / shape)
	return scale * math.Pow(-math.Log(1-u), 1.0/shape)
}

// WindTurbinePower calculates the instantaneous power output of a wind turbine in Watts (W),
// based on the standard aerodynamic power equation.
//
// This function models the power extracted from the wind using the turbine's rotor area,
// air density, wind speed, and efficiency coefficient. It's derived from the kinetic energy
// of the wind and is essential for simulating wind farm performance or energy yield predictions.
// Note: This assumes steady-state conditions and does not account for dynamic effects like
// turbulence or wake interactions.
//
// Parameters:
//   - rho: Air density in kg/m³ (typically around 1.225 at sea level and standard temperature).
//   - area: Swept area of the turbine rotor in m² (π * r² where r is the blade radius).
//   - windSpeed: Current wind speed in m/s.
//   - cp: Power coefficient of the turbine (Betz limit is 0.593, practical values range from
//     0.35 to 0.45 depending on design and operating conditions).
//
// Returns:
//
//	The power output in Watts (W).
func WindTurbinePower(rho, area, windSpeed, cp float64) float64 {
	// Formula: P = 0.5 * rho * area * windSpeed³ * cp
	return 0.5 * rho * area * math.Pow(windSpeed, 3) * cp
}
