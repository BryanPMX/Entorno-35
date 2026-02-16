package services

import (
	"errors"
	"testing"
	"time"

	"github.com/entorno35/backend/internal/core/scoring"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPublicAssessmentRepo struct {
	link          *domain.AssessmentLink
	linkErr       error
	assessment    *domain.Assessment
	assessmentErr error
	questions     []domain.Question
	questionsErr  error
}

func (s *stubPublicAssessmentRepo) GetByID(id uuid.UUID) (*domain.Assessment, error) {
	if s.assessmentErr != nil {
		return nil, s.assessmentErr
	}
	return s.assessment, nil
}

func (s *stubPublicAssessmentRepo) GetLinkByToken(token string) (*domain.AssessmentLink, error) {
	if s.linkErr != nil {
		return nil, s.linkErr
	}
	return s.link, nil
}

func (s *stubPublicAssessmentRepo) GetQuestionsByGuideTypeWithRelations(guideType domain.GuideType) ([]domain.Question, error) {
	if s.questionsErr != nil {
		return nil, s.questionsErr
	}
	return s.questions, nil
}

func (s *stubPublicAssessmentRepo) GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Assessment, error) {
	panic("unexpected call to GetByIDAndCompany")
}

func (s *stubPublicAssessmentRepo) GetResponsesByAssessmentID(assessmentID uuid.UUID) ([]domain.Response, error) {
	panic("unexpected call to GetResponsesByAssessmentID")
}

func (s *stubPublicAssessmentRepo) UpdateResult(assessment *domain.Assessment, result *scoring.AssessmentResult) error {
	panic("unexpected call to UpdateResult")
}

func (s *stubPublicAssessmentRepo) Create(assessment *domain.Assessment) error {
	panic("unexpected call to Create")
}

func (s *stubPublicAssessmentRepo) ListByCompany(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus) ([]domain.Assessment, error) {
	panic("unexpected call to ListByCompany")
}

func (s *stubPublicAssessmentRepo) ListByCompanyPaginated(companyID uuid.UUID, staffID *uuid.UUID, period *int, status *domain.AssessmentStatus, limit int, offset int) ([]domain.Assessment, int, error) {
	panic("unexpected call to ListByCompanyPaginated")
}

func (s *stubPublicAssessmentRepo) CreateLink(link *domain.AssessmentLink) error {
	panic("unexpected call to CreateLink")
}

func (s *stubPublicAssessmentRepo) UpdateLinkAccess(linkID uuid.UUID) error {
	panic("unexpected call to UpdateLinkAccess")
}

func (s *stubPublicAssessmentRepo) UpdateLinkAccessedAt(linkID uuid.UUID) error {
	panic("unexpected call to UpdateLinkAccessedAt")
}

func (s *stubPublicAssessmentRepo) GetQuestionsByGuideType(guideType domain.GuideType) ([]domain.Question, error) {
	panic("unexpected call to GetQuestionsByGuideType")
}

func (s *stubPublicAssessmentRepo) Delete(id uuid.UUID) error {
	panic("unexpected call to Delete")
}

func TestAssessmentService_GetPublicAssessment_NoQuestionsConfigured(t *testing.T) {
	assessmentID := uuid.New()

	repo := &stubPublicAssessmentRepo{
		link: &domain.AssessmentLink{
			Token:        "public-token",
			AssessmentID: &assessmentID,
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		assessment: &domain.Assessment{
			ID:        assessmentID,
			GuideType: domain.GuideTypeII,
			Status:    domain.AssessmentStatusPending,
		},
		questions: []domain.Question{},
	}

	service := NewAssessmentService(repo, nil, nil, nil, nil)

	assessment, questions, err := service.GetPublicAssessment("public-token")

	require.Error(t, err)
	assert.Nil(t, assessment)
	assert.Nil(t, questions)
	assert.Contains(t, err.Error(), "no questions configured for guide type II")
}

func TestAssessmentService_GetPublicAssessment_ReturnsAssessmentAndQuestions(t *testing.T) {
	assessmentID := uuid.New()

	repo := &stubPublicAssessmentRepo{
		link: &domain.AssessmentLink{
			Token:        "public-token",
			AssessmentID: &assessmentID,
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		},
		assessment: &domain.Assessment{
			ID:        assessmentID,
			GuideType: domain.GuideTypeII,
			Status:    domain.AssessmentStatusPending,
		},
		questions: []domain.Question{
			{ID: 1, GuideType: domain.GuideTypeII, Text: "Question 1"},
		},
	}

	service := NewAssessmentService(repo, nil, nil, nil, nil)

	assessment, questions, err := service.GetPublicAssessment("public-token")

	require.NoError(t, err)
	require.NotNil(t, assessment)
	assert.Equal(t, assessmentID, assessment.ID)
	assert.Len(t, questions, 1)
	assert.Equal(t, uint(1), questions[0].ID)
}

func TestAssessmentService_GetPublicAssessment_PropagatesLinkErrors(t *testing.T) {
	repo := &stubPublicAssessmentRepo{
		linkErr: errors.New("assessment link not found"),
	}

	service := NewAssessmentService(repo, nil, nil, nil, nil)

	assessment, questions, err := service.GetPublicAssessment("missing-token")

	require.Error(t, err)
	assert.Nil(t, assessment)
	assert.Nil(t, questions)
	assert.Contains(t, err.Error(), "invalid or expired assessment link")
}
