package domain

type Tags struct {
	Temperature float64 `json:"temperature"`
	Pressure    float64 `json:"pressure"`
}

type Payload struct {
	Timestamp string `json:"timestamp"`
	AssetID   string `json:"asset_id"`
	Tags      Tags   `json:"tags"`
}
