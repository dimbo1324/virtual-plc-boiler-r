package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayloadMarshal(t *testing.T) {
	p := Payload{
		Timestamp: "2026-02-27T10:00:00Z",
		AssetID:   "boiler_01",
		Tags: Tags{
			Temperature: 450.5,
			Pressure:    60.2,
		},
	}

	data, err := json.Marshal(p)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"timestamp":"2026-02-27T10:00:00Z","asset_id":"boiler_01","tags":{"temperature":450.5,"pressure":60.2}}`, string(data))
}

func TestPayloadUnmarshal(t *testing.T) {
	jsonStr := `{"timestamp":"2026-02-27T10:00:00Z","asset_id":"boiler_01","tags":{"temperature":450.5,"pressure":60.2}}`
	var p Payload
	err := json.Unmarshal([]byte(jsonStr), &p)
	assert.NoError(t, err)
	assert.Equal(t, "2026-02-27T10:00:00Z", p.Timestamp)
	assert.Equal(t, "boiler_01", p.AssetID)
	assert.Equal(t, 450.5, p.Tags.Temperature)
	assert.Equal(t, 60.2, p.Tags.Pressure)
}
