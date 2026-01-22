package entities

import "time"

type WorkoutType int

const (
	Estada WorkoutType = iota
	Rolo
)

type Workout struct {
	workoutType     WorkoutType
	startTime       time.Time
	distance        float64
	duration        int
	elevation       int
	avgPower        int
	normalizedPower int
	avgHeartRate    int
	maxHeartRate    int
	avgCadence      int
}

type WorkoutParams struct {
	WorkoutType  WorkoutType
	StartTime    time.Time
	Distance     float64
	Duration     int
	Elevation    int
	AvgPower     int
	AvgHeartRate int
	MaxHeartRate int
	AvgCadence   int
}

func NewWorkout(params WorkoutParams) *Workout {
	return &Workout{
		workoutType:  params.WorkoutType,
		startTime:    params.StartTime,
		distance:     params.Distance,
		duration:     params.Duration,
		elevation:    params.Elevation,
		avgPower:     params.AvgPower,
		avgHeartRate: params.AvgHeartRate,
		maxHeartRate: params.MaxHeartRate,
		avgCadence:   params.AvgCadence,
	}
}

func (w *Workout) SetNormalizedPower(np int) {
	w.normalizedPower = np
}
