package domain

// QuestionJSON represents the structure of questions in the JSON file
type QuestionJSON struct {
	Number      int    `json:"number"`
	Text        string `json:"text"`
	Type        string `json:"type"`
	Section     string `json:"section,omitempty"`      // Guide I only
	Subsection  string `json:"subsection,omitempty"`   // Guide I only
	Category    string `json:"category,omitempty"`     // Guide II/III only
	Domain      string `json:"domain,omitempty"`       // Guide II/III only
	Dimension   string `json:"dimension,omitempty"`    // Guide II/III only
	Polarity    string `json:"polarity,omitempty"`     // Guide II/III only
}

// GuideJSON represents a guide structure in the JSON file
type GuideJSON struct {
	Name           string         `json:"name"`
	Type           string         `json:"type"`
	TotalQuestions int            `json:"total_questions"`
	Questions      []QuestionJSON `json:"questions"`
}

// QuestionsDataJSON represents the root structure of nom035_questions.json
type QuestionsDataJSON struct {
	GuideI   GuideJSON `json:"guide_i"`
	GuideII  GuideJSON `json:"guide_ii"`
	GuideIII GuideJSON `json:"guide_iii"`
}

