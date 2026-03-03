package domain

type Tags struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
	Fuel        float64 `json:"fuel"`
	DrumLevel   float64 `json:"drum_level"`
	SteamFlow   float64 `json:"steam_flow"`
}

type Payload struct {
	Timestamp string `json:"timestamp"`
	AssetID   string `json:"asset_id"`
	Tags      Tags   `json:"tags"`
}
