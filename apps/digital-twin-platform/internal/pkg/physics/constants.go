package physics

// This section defines fundamental physical constants used in simulations involving
// atmospheric and solar energy calculations. These values are based on standard
// reference conditions and may be adjusted or made dynamic in future iterations
// (e.g., air density varying with temperature and pressure).
//
//  1. StandardAirDensity: Standard air density at sea level under normal conditions
//     (kg/m³). This is used in aerodynamic calculations, such as wind turbine power.
//     TODO: In future versions, this may be computed dynamically based on environmental
//     factors like temperature and atmospheric pressure.
//
//  2. SolarConstant: Average solar irradiance at the top of Earth's atmosphere (W/m²),
//     serving as a baseline for peak insolation calculations in solar energy models.
const (
	StandardAirDensity = 1.225
	SolarConstant      = 1361.0
)
