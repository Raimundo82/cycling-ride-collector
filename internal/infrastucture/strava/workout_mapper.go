package strava

import (
	"math"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/samber/lo"
)

func MapToWorkout(a *ActivityDto) *domain.Workout {
	startTime, err := time.Parse(time.RFC3339, a.StartDate)
	if err != nil {
		startTime = time.Time{}
	}

	wattsData := lo.Ternary(a.Watts != nil, a.Watts.WattsData, []int{})

	return &domain.Workout{
		ID:                     a.ID,
		WorkoutType:            lo.Ternary(a.IsTrainer, domain.Rolo, domain.Estrada),
		StartTime:              startTime,
		DurationInMin:          int(a.Duration / 60),
		DistanceInKm:           math.Trunc(a.Distance/1000*100) / 100,
		ElevationInMeters:      int(a.TotalElevationGain),
		AvgPowerInWatts:        int(a.AveragePower),
		NormalizedPowerInWatts: int(math.Round(NormalizedPower(wattsData))),
		AvgHeartRateInBpm:      int(a.AverageHeartRate),
		MaxHeartRateInBpm:      int(a.MaxHeartRate),
		AvgCadenceInRpm:        int(a.AverageCadence),
	}
}

func NormalizedPower(watts []int) float64 {
	avg30 := func(sum int) float64 {
		return float64(sum) / 30.0
	}

	if len(watts) < 30 {
		return 0
	}

	var sum int
	for i := range 30 {
		sum += watts[i]
	}

	qSum := avg30(sum) * avg30(sum) * avg30(sum) * avg30(sum)
	n := 1

	for i := 30; i < len(watts); i++ {
		sum += watts[i] - watts[i-30]
		qSum += avg30(sum) * avg30(sum) * avg30(sum) * avg30(sum)
		n++
	}

	return math.Pow(qSum/float64(n), 0.25)
}
