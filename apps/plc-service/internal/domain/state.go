package domain

import "time"

type State struct {
	Timestamp     time.Time
	FurnaceTemp   float64
	SteamPressure float64
	DrumLevel     float64
	SteamFlow     float64
}

type Controls struct {
	FuelValve      float64
	FeedwaterValve float64
	SteamValve     float64
}
