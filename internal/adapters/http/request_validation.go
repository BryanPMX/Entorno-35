package http

import (
	"fmt"
	"time"
)

const (
	defaultListLimit           = 50
	maxListLimit               = 100
	assessmentPeriodYearWindow = 1
)

func normalizeListPagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = defaultListLimit
	} else if limit > maxListLimit {
		limit = maxListLimit
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func validateAssessmentPeriod(period int, now time.Time) error {
	currentYear := now.Year()
	minYear := currentYear - assessmentPeriodYearWindow
	maxYear := currentYear + assessmentPeriodYearWindow

	if period < minYear || period > maxYear {
		return fmt.Errorf("period must be between %d and %d", minYear, maxYear)
	}

	return nil
}
