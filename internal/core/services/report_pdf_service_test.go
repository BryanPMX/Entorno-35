package services

import (
	"testing"
	"time"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportPDFService_GenerateIndividualReportPDF(t *testing.T) {
	service := NewReportPDFService()

	// Create test report data
	completedAt := time.Now()
	report := &domain.IndividualReportDTO{
		AssessmentID: uuid.New(),
		Period:       2025,
		GuideType:    domain.GuideTypeII,
		StaffName:    "Juan Pérez García",
		Department:   "Producción",
		Shift:        "Diurno",
		TotalScore:   45.5,
		RiskLevel:    domain.RiskLevelAlto,
		CategoryScores: map[string]float64{
			"Factores propios de la actividad":     25.0,
			"Ambiente de trabajo":                   8.5,
			"Organización del tiempo de trabajo":   7.0,
			"Liderazgo y relaciones en el trabajo": 5.0,
		},
		CategoryMaxScores: map[string]float64{
			"Factores propios de la actividad":     40.0,
			"Ambiente de trabajo":                   12.0,
			"Organización del tiempo de trabajo":   12.0,
			"Liderazgo y relaciones en el trabajo": 38.0,
		},
		CategoryRiskLevels: map[string]string{
			"Factores propios de la actividad":     "bajo",
			"Ambiente de trabajo":                   "medio",
			"Organización del tiempo de trabajo":   "medio",
			"Liderazgo y relaciones en el trabajo": "nulo",
		},
		DomainScores: map[string]float64{
			"Carga de trabajo":                         12.5,
			"Falta de control sobre el trabajo":        8.0,
			"Jornada de trabajo":                        5.0,
			"Liderazgo":                                10.0,
			"Relaciones en el trabajo":                 10.0,
		},
		DomainMaxScores: map[string]float64{
			"Carga de trabajo":                         24.0,
			"Falta de control sobre el trabajo":        14.0,
			"Jornada de trabajo":                        6.0,
			"Liderazgo":                                11.0,
			"Relaciones en el trabajo":                 14.0,
		},
		DomainRiskLevels: map[string]string{
			"Carga de trabajo":                         "medio",
			"Falta de control sobre el trabajo":        "bajo",
			"Jornada de trabajo":                        "muy_alto",
			"Liderazgo":                                "muy_alto",
			"Relaciones en el trabajo":                 "bajo",
		},
		RequiresMedical: false,
		CompletedAt:     &completedAt,
		Recommendations: []string{
			"Revisar y redistribuir la carga de trabajo. Implementar pausas activas.",
			"Capacitar a supervisores en liderazgo positivo.",
			"NIVEL DE RIESGO ALTO: Se recomienda implementar medidas preventivas.",
		},
	}

	// Generate PDF
	pdfData, err := service.GenerateIndividualReportPDF(report, "Test Company")
	require.NoError(t, err, "PDF generation should not return error")
	require.NotNil(t, pdfData, "PDF data should not be nil")

	// Verify PDF data
	assert.Greater(t, len(pdfData), 1000, "PDF should contain substantial data")
	
	// Verify PDF header (first 4 bytes should be %PDF)
	assert.Equal(t, byte('%'), pdfData[0], "PDF should start with %PDF header")
	assert.Equal(t, byte('P'), pdfData[1])
	assert.Equal(t, byte('D'), pdfData[2])
	assert.Equal(t, byte('F'), pdfData[3])
}

func TestReportPDFService_FormatRiskLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected string
	}{
		{"nulo", "nulo", "Nulo"},
		{"bajo", "bajo", "Bajo"},
		{"medio", "medio", "Medio"},
		{"alto", "alto", "Alto"},
		{"muy_alto", "muy_alto", "Muy Alto"},
		{"unknown", "unknown", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRiskLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReportPDFService_GetRiskColor(t *testing.T) {
	tests := []struct {
		name  string
		level string
		r     int
		g     int
		b     int
	}{
		{"nulo", "nulo", 21, 128, 61},
		{"bajo", "bajo", 77, 124, 15},
		{"medio", "medio", 161, 98, 7},
		{"alto", "alto", 194, 65, 12},
		{"muy_alto", "muy_alto", 185, 28, 28},
		{"default", "unknown", 21, 128, 61}, // defaults to nulo
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, g, b := getRiskColor(tt.level)
			assert.Equal(t, tt.r, r, "Red component should match")
			assert.Equal(t, tt.g, g, "Green component should match")
			assert.Equal(t, tt.b, b, "Blue component should match")
		})
	}
}

func TestReportPDFService_GenerateWithEmptyRecommendations(t *testing.T) {
	service := NewReportPDFService()

	report := &domain.IndividualReportDTO{
		AssessmentID:   uuid.New(),
		Period:         2025,
		GuideType:      domain.GuideTypeII,
		StaffName:      "Test User",
		Department:     "Test",
		TotalScore:     15.0,
		RiskLevel:      domain.RiskLevelBajo,
		CategoryScores: map[string]float64{
			"Test Category": 10.0,
		},
		CategoryMaxScores: map[string]float64{
			"Test Category": 40.0,
		},
		CategoryRiskLevels: map[string]string{
			"Test Category": "nulo",
		},
		DomainScores: map[string]float64{
			"Test Domain": 5.0,
		},
		DomainMaxScores: map[string]float64{
			"Test Domain": 24.0,
		},
		DomainRiskLevels: map[string]string{
			"Test Domain": "nulo",
		},
		RequiresMedical: false,
		Recommendations: []string{}, // Empty recommendations for low risk
	}

	pdfData, err := service.GenerateIndividualReportPDF(report, "Test Company")
	require.NoError(t, err, "Should generate PDF even with empty recommendations")
	require.NotNil(t, pdfData)
	assert.Greater(t, len(pdfData), 500, "PDF should still contain data")
}
