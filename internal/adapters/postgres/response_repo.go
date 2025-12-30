package postgres

import (
	"errors"
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// responseRepository implements the ResponseRepository interface using GORM
type responseRepository struct {
	db *gorm.DB
}

// NewResponseRepository creates a new Postgres implementation of ResponseRepository
func NewResponseRepository(db *gorm.DB) ports.ResponseRepository {
	return &responseRepository{db: db}
}

// SaveResponses saves multiple responses in a transaction
// If any response fails to save, the entire transaction is rolled back
func (r *responseRepository) SaveResponses(responses []domain.Response) error {
	if len(responses) == 0 {
		return fmt.Errorf("no responses to save")
	}

	// Start a transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i := range responses {
			response := &responses[i]

			// Ensure AssessmentID is set (critical constraint)
			if response.AssessmentID == uuid.Nil {
				return fmt.Errorf("assessment ID is required for all responses")
			}

			// Ensure QuestionID is set
			if response.QuestionID == 0 {
				return fmt.Errorf("question ID is required for all responses")
			}

			// Validate SelectedValue range (0-4)
			if response.SelectedValue < 0 || response.SelectedValue > 4 {
				return fmt.Errorf("selected value must be between 0 and 4, got %d for question %d", response.SelectedValue, response.QuestionID)
			}

			// Upsert response (create or update)
			// Note: Based on schema, there should be a unique constraint on (assessment_id, question_id)
			// First, try to find existing response
			var existingResponse domain.Response
			result := tx.Where("assessment_id = ? AND question_id = ?", response.AssessmentID, response.QuestionID).
				First(&existingResponse)

			if result.Error == nil {
				// Update existing response
				existingResponse.SelectedValue = response.SelectedValue
				existingResponse.CalculatedScore = response.CalculatedScore
				existingResponse.AnsweredAt = response.AnsweredAt

				if err := tx.Save(&existingResponse).Error; err != nil {
					return fmt.Errorf("failed to update response for question %d: %w", response.QuestionID, err)
				}
			} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Create new response
				if err := tx.Create(response).Error; err != nil {
					return fmt.Errorf("failed to create response for question %d: %w", response.QuestionID, err)
				}
			} else {
				return fmt.Errorf("failed to check existing response for question %d: %w", response.QuestionID, result.Error)
			}
		}

		return nil
	})
}

