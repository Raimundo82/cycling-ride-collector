package activity_excel

import (
	"fmt"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/xuri/excelize/v2"
)

const defaultSheet = "Sheet1"

type excelWorkoutPeriodSaver struct {
	filePath     string
	templatePath string
	sheetName    string
	startCell    string
}

// SaveAll implements [contracts.WorkoutRepository].
func (e *excelWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout, athlete *domain.Athlete) error {
	f, err := e.openOrCreate()
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	startCol, startRow, err := excelize.CellNameToCoordinates(e.startCell)
	if err != nil {
		return fmt.Errorf("invalid start cell %q: %w", e.startCell, err)
	}

	for i, w := range workouts {
		row := workoutToRow(w, athlete.WeightInKg())
		for j, val := range row {
			cell, cellErr := excelize.CoordinatesToCellName(startCol+j, startRow+i)
			if cellErr != nil {
				return fmt.Errorf("coordinate error at row %d col %d: %w", i, j, cellErr)
			}
			if err := f.SetCellValue(e.sheetName, cell, val); err != nil {
				return fmt.Errorf("failed to set cell %s: %w", cell, err)
			}
		}
	}

	return f.SaveAs(e.filePath)
}

func (e *excelWorkoutPeriodSaver) openOrCreate() (*excelize.File, error) {
	if e.templatePath != "" {
		f, err := excelize.OpenFile(e.templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open template %q: %w", e.templatePath, err)
		}
		return f, nil
	}
	return excelize.NewFile(), nil
}

func NewExcelWorkoutPeriodSaver(filePath string) contracts.WorkoutRepository {
	return &excelWorkoutPeriodSaver{
		filePath:  filePath,
		sheetName: defaultSheet,
		startCell: "A1",
	}
}

func NewExcelWorkoutPeriodSaverWithOptions(filePath, templatePath, sheetName, startCell string) contracts.WorkoutRepository {
	return &excelWorkoutPeriodSaver{
		filePath:     filePath,
		templatePath: templatePath,
		sheetName:    sheetName,
		startCell:    startCell,
	}
}

var _ contracts.WorkoutRepository = (*excelWorkoutPeriodSaver)(nil)
