package domain

type Athlete struct {
	weightInKg              float64
	heartRateThresholdInBpm int
	powerThresholdInWatts   int
}

func NewAthlete(weightInKg float64, heartRateThresholdInBpm, powerThresholdInWatts int) *Athlete {
	a := &Athlete{}
	a.setWeightInKg(weightInKg)
	a.setHeartRateThresholdInBpm(heartRateThresholdInBpm)
	a.setPowerThresholdInWatts(powerThresholdInWatts)
	return a
}

func (a *Athlete) WeightInKg() float64 {
	return a.weightInKg
}

func (a *Athlete) HeartRateThresholdInBpm() int {
	return a.heartRateThresholdInBpm
}

func (a *Athlete) PowerThresholdInWatts() int {
	return a.powerThresholdInWatts
}

func (a *Athlete) setWeightInKg(weightInKg float64) {
	if weightInKg < 0 {
		weightInKg = 0
	}
	a.weightInKg = weightInKg
}

func (a *Athlete) setHeartRateThresholdInBpm(heartRateThresholdInBpm int) {
	if heartRateThresholdInBpm < 0 {
		heartRateThresholdInBpm = 0
	}
	a.heartRateThresholdInBpm = heartRateThresholdInBpm
}

func (a *Athlete) setPowerThresholdInWatts(powerThresholdInWatts int) {
	if powerThresholdInWatts < 0 {
		powerThresholdInWatts = 0
	}
	a.powerThresholdInWatts = powerThresholdInWatts
}
