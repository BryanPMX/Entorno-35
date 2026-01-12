package services

import (
	"fmt"

	"github.com/entorno35/backend/internal/core/ports"
	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
)

// ReportService handles report generation operations
type ReportService struct {
	reportRepo ports.ReportRepository
}

// NewReportService creates a new report service
func NewReportService(reportRepo ports.ReportRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
	}
}

// GenerateIndividualReport generates an individual assessment report with recommendations
func (s *ReportService) GenerateIndividualReport(assessmentID uuid.UUID, companyID uuid.UUID) (*domain.IndividualReportDTO, error) {
	// Fetch report data from repository (includes dynamic max scores calculation)
	report, err := s.reportRepo.GetIndividualReport(assessmentID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get individual report: %w", err)
	}

	// Generate recommendations based on risk level and domain scores
	report.Recommendations = s.generateRecommendations(report.RiskLevel, report.DomainScores, report.GuideType)

	return report, nil
}

// GenerateGeneralReport generates a company-wide general report
func (s *ReportService) GenerateGeneralReport(companyID uuid.UUID, period *int) (*domain.GeneralReportDTO, error) {
	report, err := s.reportRepo.GetGeneralReport(companyID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get general report: %w", err)
	}

	return report, nil
}

// generateRecommendations creates recommendations based on NOM-035 Section 8 guidelines
// Returns recommendations for high-risk domains when risk level is Alto or Muy Alto
func (s *ReportService) generateRecommendations(riskLevel domain.RiskLevel, domainScores map[string]float64, guideType domain.GuideType) []string {
	recommendations := []string{}

	// Only generate recommendations for Alto or Muy Alto risk levels
	if riskLevel != domain.RiskLevelAlto && riskLevel != domain.RiskLevelMuyAlto {
		return recommendations
	}

	// NOM-035 Section 8 recommendations based on domain
	// Map domain names to specific recommendations
	domainRecommendations := map[string]string{
		"Carga de trabajo":                    "Revisar y redistribuir la carga de trabajo. Implementar pausas activas y rotación de tareas.",
		"Falta de control sobre el trabajo":   "Aumentar la autonomía y participación en la toma de decisiones. Establecer canales de comunicación efectivos.",
		"Jornada de trabajo":                   "Revisar horarios y turnos. Evitar horas extras excesivas. Implementar descansos adecuados.",
		"Interferencia en la relación trabajo-familia": "Establecer políticas de conciliación trabajo-familia. Flexibilizar horarios cuando sea posible.",
		"Liderazgo":                           "Capacitar a supervisores en liderazgo positivo. Establecer retroalimentación constructiva.",
		"Relaciones en el trabajo":           "Fomentar el trabajo en equipo. Implementar programas de integración y comunicación.",
		"Violencia":                           "Implementar protocolos de prevención y atención de violencia laboral. Capacitar al personal.",
		"Reconocimiento del desempeño":        "Establecer sistemas de reconocimiento y evaluación justos. Comunicar expectativas claras.",
		"Insuficiente sentido de pertenencia": "Fortalecer la cultura organizacional. Mejorar la comunicación interna y participación.",
		"Entorno organizacional":             "Revisar políticas organizacionales. Mejorar condiciones de trabajo y ambiente laboral.",
	}

	// Generate recommendations for domains with high scores
	// For Guide II/III, we check domain scores
	// Threshold: if domain score is above average or in top domains, recommend action
	if guideType == domain.GuideTypeII || guideType == domain.GuideTypeIII {
		// Find domains with highest scores (top 3)
		topDomains := s.getTopDomains(domainScores, 3)
		
		for _, domainName := range topDomains {
			if rec, exists := domainRecommendations[domainName]; exists {
				recommendations = append(recommendations, rec)
			}
		}
	}

	// Add general recommendation based on risk level
	if riskLevel == domain.RiskLevelMuyAlto {
		recommendations = append(recommendations, "NIVEL DE RIESGO MUY ALTO: Se requiere intervención inmediata. Implementar medidas correctivas urgentes y seguimiento médico especializado.")
	} else if riskLevel == domain.RiskLevelAlto {
		recommendations = append(recommendations, "NIVEL DE RIESGO ALTO: Se recomienda implementar medidas preventivas y correctivas. Realizar seguimiento periódico.")
	}

	return recommendations
}

// getTopDomains returns the top N domains by score
func (s *ReportService) getTopDomains(domainScores map[string]float64, n int) []string {
	if len(domainScores) == 0 {
		return []string{}
	}

	// Create slice of domain-score pairs
	type domainScore struct {
		name  string
		score float64
	}

	domains := make([]domainScore, 0, len(domainScores))
	for name, score := range domainScores {
		domains = append(domains, domainScore{name: name, score: score})
	}

	// Sort by score (descending) - simple bubble sort for small datasets
	for i := 0; i < len(domains)-1; i++ {
		for j := i + 1; j < len(domains); j++ {
			if domains[j].score > domains[i].score {
				domains[i], domains[j] = domains[j], domains[i]
			}
		}
	}

	// Return top N
	topN := n
	if topN > len(domains) {
		topN = len(domains)
	}

	result := make([]string, 0, topN)
	for i := 0; i < topN; i++ {
		result = append(result, domains[i].name)
	}

	return result
}


