package domain

import (
	"time"
)

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
	legSensations          LegSensations
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
	LegSensations          string
}

func NewWorkout(params WorkoutParams) *Workout {
	workout := &Workout{
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
	workout.SetLegSensations(params.LegSensations)
	return workout
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

func (w *Workout) SetLegSensations(ls string) {
	switch ls {
	case "Muito Más":
		w.legSensations = VeryBad
	case "Más":
		w.legSensations = Bad
	case "Médias":
		w.legSensations = Medium
	case "Boas":
		w.legSensations = Good
	case "Muito Boas":
		w.legSensations = VeryGood
	case "Excelentes":
		w.legSensations = Excellent
	default:
		w.legSensations = ""
	}
}

func (w *Workout) LegSensations() LegSensations {
	return w.legSensations
}

func (w *Workout) IsRestDay() bool {
	return w.WorkoutType == Descanso
}
