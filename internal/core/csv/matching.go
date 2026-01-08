package csv

import (
	"math"
	"strings"
	"unicode"
)

// FuzzyMatcher provides fuzzy string matching capabilities
type FuzzyMatcher struct{}

// NewFuzzyMatcher creates a new fuzzy matcher
func NewFuzzyMatcher() *FuzzyMatcher {
	return &FuzzyMatcher{}
}

// FindBestMatch finds the best matching expected field for a given CSV header
func (fm *FuzzyMatcher) FindBestMatch(header string, expectedFields []ExpectedField) (*ColumnMapping, bool) {
	header = normalizeHeader(header)
	bestMatch := &ColumnMapping{
		CSVHeader:     header,
		Confidence:    0.0,
		MatchType:     MatchTypeFuzzy,
	}

	found := false

	for _, field := range expectedFields {
		// Check exact synonym match first
		for _, synonym := range field.Synonyms {
			if normalizeHeader(synonym) == header {
				return &ColumnMapping{
					CSVHeader:     header,
					ExpectedField: field.Name,
					Confidence:    1.0,
					MatchType:     MatchTypeSynonym,
				}, true
			}
		}

		// Check exact field name match
		if normalizeHeader(field.Name) == header {
			return &ColumnMapping{
				CSVHeader:     header,
				ExpectedField: field.Name,
				Confidence:    1.0,
				MatchType:     MatchTypeExact,
			}, true
		}

		// Calculate fuzzy similarity for all synonyms
		maxSimilarity := 0.0
		for _, synonym := range field.Synonyms {
			similarity := fm.calculateSimilarity(header, normalizeHeader(synonym))
			if similarity > maxSimilarity {
				maxSimilarity = similarity
			}
		}

		// Also check similarity with field name itself
		fieldSimilarity := fm.calculateSimilarity(header, normalizeHeader(field.Name))
		if fieldSimilarity > maxSimilarity {
			maxSimilarity = fieldSimilarity
		}

		// Update best match if this field has higher similarity
		if maxSimilarity > bestMatch.Confidence {
			bestMatch.ExpectedField = field.Name
			bestMatch.Confidence = maxSimilarity
			bestMatch.MatchType = MatchTypeFuzzy
			found = true
		}
	}

	if found && bestMatch.Confidence >= 0.6 { // Minimum confidence threshold
		return bestMatch, true
	}

	return nil, false
}

// calculateSimilarity calculates string similarity using multiple algorithms
func (fm *FuzzyMatcher) calculateSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}

	// Use a combination of algorithms for better accuracy
	levenshteinSim := 1.0 - float64(fm.levenshteinDistance(a, b))/math.Max(float64(len(a)), float64(len(b)))
	jaccardSim := fm.jaccardSimilarity(a, b)

	// Weighted average: give more weight to Levenshtein for short strings
	if len(a) <= 10 || len(b) <= 10 {
		return levenshteinSim*0.7 + jaccardSim*0.3
	}

	return levenshteinSim*0.5 + jaccardSim*0.5
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func (fm *FuzzyMatcher) levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			matrix[i][j] = int(math.Min(
				float64(matrix[i-1][j]+1),      // deletion
				math.Min(
					float64(matrix[i][j-1]+1),  // insertion
					float64(matrix[i-1][j-1]+cost), // substitution
				),
			))
		}
	}

	return matrix[len(a)][len(b)]
}

// jaccardSimilarity calculates Jaccard similarity between two strings
func (fm *FuzzyMatcher) jaccardSimilarity(a, b string) float64 {
	// Convert strings to sets of bigrams
	setA := fm.getBigrams(a)
	setB := fm.getBigrams(b)

	if len(setA) == 0 && len(setB) == 0 {
		return 1.0
	}

	intersection := 0
	setAMap := make(map[string]bool)
	for _, bigram := range setA {
		setAMap[bigram] = true
	}

	for _, bigram := range setB {
		if setAMap[bigram] {
			intersection++
		}
	}

	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// getBigrams returns all bigrams (pairs of consecutive characters) from a string
func (fm *FuzzyMatcher) getBigrams(s string) []string {
	s = strings.ToLower(s)
	var bigrams []string

	for i := 0; i < len(s)-1; i++ {
		bigrams = append(bigrams, s[i:i+2])
	}

	return bigrams
}

// normalizeHeader normalizes a CSV header for comparison
func normalizeHeader(header string) string {
	// Convert to lowercase
	header = strings.ToLower(header)

	// Remove accents and special characters
	header = removeAccents(header)

	// Replace common separators with underscores
	header = strings.ReplaceAll(header, " ", "_")
	header = strings.ReplaceAll(header, "-", "_")
	header = strings.ReplaceAll(header, ".", "_")

	// Remove non-alphanumeric characters except underscores
	var result strings.Builder
	for _, r := range header {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// removeAccents removes Spanish accents from characters
func removeAccents(s string) string {
	accents := map[rune]rune{
		'á': 'a', 'é': 'e', 'í': 'i', 'ó': 'o', 'ú': 'u',
		'Á': 'a', 'É': 'e', 'Í': 'i', 'Ó': 'o', 'Ú': 'u',
		'ñ': 'n', 'Ñ': 'n',
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, exists := accents[r]; exists {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}