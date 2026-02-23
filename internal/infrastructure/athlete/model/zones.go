package athlete_model

type Zone struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type PowerRangeZones struct {
	Zones []Zone `json:"zones"`
}

type HeartRateRangeZones struct {
	Zones []Zone `json:"zones"`
}

type Zones struct {
	PowerRangeZones     *PowerRangeZones     `json:"power"`
	HeartRateRangeZones *HeartRateRangeZones `json:"heart_rate"`
}
