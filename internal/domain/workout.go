package domain

import "time"

type WorkoutType int

const (
	Descanso WorkoutType = iota
	Estrada
	Rolo
	Prova
)

type LegSensations string

const (
	VeryBad   LegSensations = "Muito Más"
	Bad       LegSensations = "Más"
	Medium    LegSensations = "Médias"
	Good      LegSensations = "Boas"
	VeryGood  LegSensations = "Muito Boas"
	Excellent LegSensations = "Excelentes"
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
	LegSensations          LegSensations
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
	LegSensations          LegSensations
}

func NewWorkout(params *WorkoutParams) *Workout {
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
		LegSensations:          params.LegSensations,
	}
}

func (wt WorkoutType) String() string {
	switch wt {
	case Estrada:
		return "Estrada"
	case Rolo:
		return "Rolo"
	case Prova:
		return "Prova"
	default:
		return "Descanso"
	}
}
