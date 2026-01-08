package csv

import (
	"regexp"
	"strconv"
	"strings"
)

// DataTypeInferer analyzes sample data to infer column types
type DataTypeInferer struct{}

// NewDataTypeInferer creates a new data type inferer
func NewDataTypeInferer() *DataTypeInferer {
	return &DataTypeInferer{}
}

// InferType analyzes a sample of values from a CSV column and infers the data type
func (dti *DataTypeInferer) InferType(sampleValues []string, header string) (FieldType, float64) {
	if len(sampleValues) == 0 {
		return FieldTypeText, 0.5
	}

	// Remove empty values for analysis
	var nonEmptyValues []string
	for _, v := range sampleValues {
		if strings.TrimSpace(v) != "" {
			nonEmptyValues = append(nonEmptyValues, v)
		}
	}

	if len(nonEmptyValues) == 0 {
		return FieldTypeText, 0.5
	}

	// Check confidence for each type
	emailConfidence := dti.checkEmailConfidence(nonEmptyValues)
	numberConfidence := dti.checkNumberConfidence(nonEmptyValues)
	textConfidence := dti.checkTextConfidence(nonEmptyValues, header)
	enumConfidence := dti.checkEnumConfidence(nonEmptyValues, header)

	// Return the type with highest confidence
	typeConfidences := map[FieldType]float64{
		FieldTypeEmail: emailConfidence,
		FieldTypeNumber: numberConfidence,
		FieldTypeText:   textConfidence,
		FieldTypeEnum:   enumConfidence,
	}

	maxConfidence := 0.0
	bestType := FieldTypeText

	for fieldType, confidence := range typeConfidences {
		if confidence > maxConfidence {
			maxConfidence = confidence
			bestType = fieldType
		}
	}

	return bestType, maxConfidence
}

// checkEmailConfidence determines how likely a column contains email addresses
func (dti *DataTypeInferer) checkEmailConfidence(values []string) float64 {
	if len(values) == 0 {
		return 0.0
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	emailCount := 0

	for _, value := range values {
		if emailRegex.MatchString(strings.TrimSpace(value)) {
			emailCount++
		}
	}

	return float64(emailCount) / float64(len(values))
}

// checkNumberConfidence determines how likely a column contains numbers
func (dti *DataTypeInferer) checkNumberConfidence(values []string) float64 {
	if len(values) == 0 {
		return 0.0
	}

	numberCount := 0
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		// Check if it's a number
		if _, err := strconv.ParseFloat(trimmed, 64); err == nil {
			numberCount++
		} else if len(trimmed) == 18 && dti.looksLikeCURP(trimmed) {
			// CURP is 18 characters and looks like an ID
			numberCount++
		}
	}

	return float64(numberCount) / float64(len(values))
}

// checkTextConfidence determines how likely a column contains text
func (dti *DataTypeInferer) checkTextConfidence(values []string, header string) float64 {
	if len(values) == 0 {
		return 0.5
	}

	// Text is the default type, so we give it a base confidence
	baseConfidence := 0.7

	// Increase confidence if header suggests text content
	header = strings.ToLower(header)
	textHeaders := []string{"name", "full", "employee", "staff", "area", "department", "job", "position", "role"}

	for _, textHeader := range textHeaders {
		if strings.Contains(header, textHeader) {
			baseConfidence += 0.2
			break
		}
	}

	// Decrease confidence if most values look like numbers or emails
	emailConfidence := dti.checkEmailConfidence(values)
	numberConfidence := dti.checkNumberConfidence(values)

	// Reduce text confidence if other types are more likely
	if emailConfidence > 0.8 {
		baseConfidence -= 0.3
	}
	if numberConfidence > 0.8 {
		baseConfidence -= 0.3
	}

	if baseConfidence < 0.1 {
		baseConfidence = 0.1
	}
	if baseConfidence > 1.0 {
		baseConfidence = 1.0
	}

	return baseConfidence
}

// checkEnumConfidence determines how likely a column contains enum values
func (dti *DataTypeInferer) checkEnumConfidence(values []string, header string) float64 {
	if len(values) == 0 {
		return 0.0
	}

	// Check if values look like common enums
	header = strings.ToLower(header)
	uniqueValues := dti.getUniqueValues(values)

	// Gender-like fields
	if strings.Contains(header, "gender") || strings.Contains(header, "sex") ||
	   strings.Contains(header, "genero") || strings.Contains(header, "sexo") {
		genderValues := []string{"masculino", "femenino", "male", "female", "m", "f", "otro", "other"}
		if dti.matchesEnumValues(uniqueValues, genderValues) {
			return 0.9
		}
	}

	// Shift-like fields
	if strings.Contains(header, "shift") || strings.Contains(header, "turno") ||
	   strings.Contains(header, "jornada") {
		shiftValues := []string{"diurno", "nocturno", "mixto", "day", "night", "mixed", "regular"}
		if dti.matchesEnumValues(uniqueValues, shiftValues) {
			return 0.9
		}
	}

	// Low uniqueness suggests enum (few distinct values relative to total)
	uniquenessRatio := float64(len(uniqueValues)) / float64(len(values))
	if uniquenessRatio < 0.3 && len(uniqueValues) <= 10 {
		return 0.7
	}

	return 0.0
}

// matchesEnumValues checks if unique values match known enum patterns
func (dti *DataTypeInferer) matchesEnumValues(uniqueValues []string, enumValues []string) bool {
	matches := 0
	for _, uniqueVal := range uniqueValues {
		uniqueLower := strings.ToLower(strings.TrimSpace(uniqueVal))
		for _, enumVal := range enumValues {
			if uniqueLower == strings.ToLower(enumVal) {
				matches++
				break
			}
		}
	}

	// At least 70% of unique values should match known enum values
	return float64(matches)/float64(len(uniqueValues)) >= 0.7
}

// getUniqueValues returns unique non-empty values from a slice
func (dti *DataTypeInferer) getUniqueValues(values []string) []string {
	seen := make(map[string]bool)
	var unique []string

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			unique = append(unique, trimmed)
		}
	}

	return unique
}

// looksLikeCURP checks if a string resembles a CURP format
func (dti *DataTypeInferer) looksLikeCURP(value string) bool {
	if len(value) != 18 {
		return false
	}

	// CURP format: 4 letters, 6 digits, 6 letters, 2 digits
	// This is a basic check - real validation is more complex
	letters := 0
	digits := 0

	for _, r := range value {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			letters++
		} else if r >= '0' && r <= '9' {
			digits++
		}
	}

	// CURP has 10 letters and 8 digits (approximately)
	return letters >= 8 && digits >= 6
}