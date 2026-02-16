package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeListPagination(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		limit          int
		offset         int
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "uses default limit when limit is zero",
			limit:          0,
			offset:         15,
			expectedLimit:  defaultListLimit,
			expectedOffset: 15,
		},
		{
			name:           "uses default limit when limit is negative",
			limit:          -10,
			offset:         -5,
			expectedLimit:  defaultListLimit,
			expectedOffset: 0,
		},
		{
			name:           "caps limit at max",
			limit:          101,
			offset:         8,
			expectedLimit:  maxListLimit,
			expectedOffset: 8,
		},
		{
			name:           "keeps values when valid",
			limit:          25,
			offset:         0,
			expectedLimit:  25,
			expectedOffset: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			limit, offset := normalizeListPagination(tc.limit, tc.offset)
			assert.Equal(t, tc.expectedLimit, limit)
			assert.Equal(t, tc.expectedOffset, offset)
		})
	}
}

func TestValidateAssessmentPeriod(t *testing.T) {
	t.Parallel()

	fixedNow := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	assert.NoError(t, validateAssessmentPeriod(2025, fixedNow))
	assert.NoError(t, validateAssessmentPeriod(2026, fixedNow))
	assert.NoError(t, validateAssessmentPeriod(2027, fixedNow))
	assert.EqualError(t, validateAssessmentPeriod(2024, fixedNow), "period must be between 2025 and 2027")
	assert.EqualError(t, validateAssessmentPeriod(2028, fixedNow), "period must be between 2025 and 2027")
}
