package domain

type WorkoutType int

const (
	Estrada WorkoutType = iota
	Rolo
)

type Workout struct {
	Id              int64
	WorkoutType     WorkoutType
	StartTime       string
	Distance        float64
	Duration        int
	Elevation       int
	AvgPower        int
	NormalizedPower int
	AvgHeartRate    int
	MaxHeartRate    int
	AvgCadence      int
}

type WorkoutParams struct {
	Id           int64
	WorkoutType  WorkoutType
	StartTime    string
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
		Id:           params.Id,
		WorkoutType:  params.WorkoutType,
		StartTime:    params.StartTime,
		Distance:     params.Distance,
		Duration:     params.Duration,
		Elevation:    params.Elevation,
		AvgPower:     params.AvgPower,
		AvgHeartRate: params.AvgHeartRate,
		MaxHeartRate: params.MaxHeartRate,
		AvgCadence:   params.AvgCadence,
	}
}
