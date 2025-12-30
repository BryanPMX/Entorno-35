package scoring

// RiskLevel represents the risk level classification
type RiskLevel string

const (
	RiskLevelNulo     RiskLevel = "nulo"
	RiskLevelBajo     RiskLevel = "bajo"
	RiskLevelMedio    RiskLevel = "medio"
	RiskLevelAlto     RiskLevel = "alto"
	RiskLevelMuyAlto  RiskLevel = "muy_alto"
)

// RiskThresholds defines the score ranges for risk levels
type RiskThresholds struct {
	Ranges [4]float64 // [Limit1, Limit2, Limit3, Limit4] - 4 thresholds define 5 levels
}

// GetRiskLevel calculates the risk level based on score and thresholds
// Logic from official NOM-035 document:
//   Nulo: Score < Limit1
//   Bajo: Limit1 <= Score < Limit2
//   Medio: Limit2 <= Score < Limit3
//   Alto: Limit3 <= Score < Limit4
//   Muy Alto: Score >= Limit4
func (rt *RiskThresholds) GetRiskLevel(score float64) RiskLevel {
	if score < rt.Ranges[0] {
		return RiskLevelNulo
	}
	if score < rt.Ranges[1] {
		return RiskLevelBajo
	}
	if score < rt.Ranges[2] {
		return RiskLevelMedio
	}
	if score < rt.Ranges[3] {
		return RiskLevelAlto
	}
	return RiskLevelMuyAlto
}

// GuideIIScoringRules contains scoring thresholds for Guide II (16-50 employees)
type GuideIIScoringRules struct {
	TotalScore  RiskThresholds
	Categories  map[string]RiskThresholds
	Domains     map[string]RiskThresholds
}

// GuideIIIScoringRules contains scoring thresholds for Guide III (>50 employees)
type GuideIIIScoringRules struct {
	TotalScore  RiskThresholds
	Categories  map[string]RiskThresholds
	Domains     map[string]RiskThresholds
}

// ScoringRules contains all scoring rules for all guides
type ScoringRules struct {
	GuideII  GuideIIScoringRules
	GuideIII GuideIIIScoringRules
}

// LoadScoringRules returns the hardcoded scoring rules from NOM-035 official document
// These are static business rules that don't change, so they're configured here
// Source: Official NOM-035-STPS-2018 document (Pages 28 and 36)
func LoadScoringRules() *ScoringRules {
	return &ScoringRules{
		GuideII: GuideIIScoringRules{
			TotalScore: RiskThresholds{Ranges: [4]float64{20, 45, 70, 90}},
			Categories: map[string]RiskThresholds{
				"Ambiente de trabajo":                        {Ranges: [4]float64{3, 5, 7, 9}},
				"Factores propios de la actividad":          {Ranges: [4]float64{10, 20, 30, 40}},
				"Organización del tiempo de trabajo":        {Ranges: [4]float64{4, 6, 9, 12}},
				"Liderazgo y relaciones en el trabajo":      {Ranges: [4]float64{10, 18, 28, 38}},
			},
			Domains: map[string]RiskThresholds{
				"Condiciones en el ambiente de trabajo":    {Ranges: [4]float64{3, 5, 7, 9}},
				"Carga de trabajo":                         {Ranges: [4]float64{12, 16, 20, 24}},
				"Falta de control sobre el trabajo":        {Ranges: [4]float64{5, 8, 11, 14}},
				"Jornada de trabajo":                       {Ranges: [4]float64{1, 2, 4, 6}},
				"Interferencia en la relación trabajo-familia": {Ranges: [4]float64{1, 2, 4, 6}},
				"Liderazgo":                                {Ranges: [4]float64{3, 5, 8, 11}},
				"Relaciones en el trabajo":                 {Ranges: [4]float64{5, 8, 11, 14}},
				"Violencia":                                {Ranges: [4]float64{7, 10, 13, 16}},
			},
		},
		GuideIII: GuideIIIScoringRules{
			TotalScore: RiskThresholds{Ranges: [4]float64{50, 75, 99, 140}},
			Categories: map[string]RiskThresholds{
				"Ambiente de trabajo":                        {Ranges: [4]float64{5, 9, 11, 14}},
				"Factores propios de la actividad":          {Ranges: [4]float64{15, 30, 45, 60}},
				"Organización del tiempo de trabajo":        {Ranges: [4]float64{5, 7, 10, 13}},
				"Liderazgo y relaciones en el trabajo":      {Ranges: [4]float64{14, 29, 42, 58}},
				"Entorno organizacional":                    {Ranges: [4]float64{10, 14, 18, 23}},
			},
			Domains: map[string]RiskThresholds{
				"Condiciones en el ambiente de trabajo":    {Ranges: [4]float64{5, 9, 11, 14}},
				"Carga de trabajo":                         {Ranges: [4]float64{15, 21, 27, 37}},
				"Falta de control sobre el trabajo":        {Ranges: [4]float64{11, 16, 21, 25}},
				"Jornada de trabajo":                       {Ranges: [4]float64{1, 2, 4, 6}},
				"Interferencia en la relación trabajo-familia": {Ranges: [4]float64{4, 6, 8, 10}},
				"Liderazgo":                                {Ranges: [4]float64{9, 12, 16, 20}},
				"Relaciones en el trabajo":                 {Ranges: [4]float64{10, 13, 17, 21}},
				"Violencia":                                {Ranges: [4]float64{7, 10, 13, 16}},
				"Reconocimiento del desempeño":             {Ranges: [4]float64{6, 10, 14, 18}},
				"Insuficiente sentido de pertenencia e inestabilidad": {Ranges: [4]float64{4, 6, 8, 10}},
			},
		},
	}
}

