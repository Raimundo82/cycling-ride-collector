package domain

import "time"

type WorkoutType int

const (
	Estrada WorkoutType = iota
	Rolo
	Mixed
)

type Workout struct {
	ID                     int64
	WorkoutType            WorkoutType
	StartTime              time.Time
	DistanceInKm           float64
	DurationInMin          int
	ElevationInMeters      int
	AvgPowerInWatts        int
	NormalizedPowerInWatts int
	AvgHeartRateInBpm      int
	MaxHeartRateInBpm      int
	AvgCadenceInRpm        int
}

type WorkoutParams struct {
	ID                     int64
	WorkoutType            WorkoutType
	StartTime              time.Time
	DistanceInKm           float64
	DurationInMin          int
	ElevationInMeters      int
	AvgPowerInWatts        int
	NormalizedPowerInWatts int
	AvgHeartRateInBpm      int
	MaxHeartRateInBpm      int
	AvgCadenceInRpm        int
}

func NewWorkout(params WorkoutParams) *Workout {
	return &Workout{
		ID:                     params.ID,
		WorkoutType:            params.WorkoutType,
		StartTime:              params.StartTime,
		DistanceInKm:           params.DistanceInKm,
		DurationInMin:          params.DurationInMin,
		ElevationInMeters:      params.ElevationInMeters,
		AvgPowerInWatts:        params.AvgPowerInWatts,
		NormalizedPowerInWatts: params.NormalizedPowerInWatts,
		MaxHeartRateInBpm:      params.MaxHeartRateInBpm,
		AvgHeartRateInBpm:      params.AvgHeartRateInBpm,
		AvgCadenceInRpm:        params.AvgCadenceInRpm,
	}
}

func (wt WorkoutType) String() string {
	switch wt {
	case Estrada:
		return "Estrada"
	case Rolo:
		return "Rolo"
	case Mixed:
		return "Mixed"
	default:
		return ""
	}
}
