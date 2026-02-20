package activity_provider

import (
	"math"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
	"github.com/samber/lo"
)

func MapToWorkout(a *activity_model.ActivityDto) *domain.Workout {
	startTime, err := time.Parse(time.RFC3339, a.StartDate)
	if err != nil {
		startTime = time.Time{}
	}

	wattsData := lo.Ternary(a.Watts != nil, a.Watts.WattsData, []int{})

	return domain.NewWorkout(&domain.WorkoutParams{
		ID:                     a.ID,
		WorkoutType:            SetWorkoutType(a),
		StartTime:              startTime,
		DurationInMin:          int(a.Duration / 60),
		DistanceInKm:           math.Trunc(a.Distance/1000*100) / 100,
		ElevationInMeters:      int(a.TotalElevationGain),
		AvgPowerInWatts:        lo.Ternary(a.DeviceWatts, int(a.AveragePower), -1),
		NormalizedPowerInWatts: lo.Ternary(a.DeviceWatts, int(math.Round(NormalizedPower(wattsData))), -1),
		AvgHeartRateInBpm:      lo.Ternary(a.HasHeartRate, int(a.AverageHeartRate), -1),
		MaxHeartRateInBpm:      lo.Ternary(a.HasHeartRate, int(a.MaxHeartRate), -1),
		AvgCadenceInRpm:        lo.Ternary(a.AverageCadence > 0, int(a.AverageCadence), -1),
		LegSensations:          lo.Ternary(a.LegSensations == "", string(domain.Medium), a.LegSensations),
	})
}

func SetWorkoutType(a *activity_model.ActivityDto) domain.WorkoutType {
	if a.IsTrainer {
		return domain.Rolo
	}

	switch a.WorkoutType {
	case 10, 12:
		return domain.Estrada
	case 11:
		return domain.Prova
	default:
		return domain.Estrada
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
