package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// DetectionService provides intelligent CSV column detection and mapping
type DetectionService struct {
	fuzzyMatcher    *FuzzyMatcher
	dataTypeInferer *DataTypeInferer
}

// NewDetectionService creates a new CSV detection service
func NewDetectionService() *DetectionService {
	return &DetectionService{
		fuzzyMatcher:    NewFuzzyMatcher(),
		dataTypeInferer: NewDataTypeInferer(),
	}
}

// AnalyzeCSV analyzes a CSV file and returns column mappings
func (ds *DetectionService) AnalyzeCSV(reader io.Reader) (*DetectionResult, error) {
	csvReader := csv.NewReader(reader)

	// Read headers
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Sample first few rows for data type inference
	sampleRows := ds.sampleRows(csvReader, 10) // Sample up to 10 rows

	return ds.DetectColumnMappings(headers, sampleRows), nil
}

// DetectColumnMappings analyzes headers and sample data to create column mappings
func (ds *DetectionService) DetectColumnMappings(headers []string, sampleData [][]string) *DetectionResult {
	result := &DetectionResult{
		Mappings:       []ColumnMapping{},
		Unmapped:       []string{},
		Confidence:     0.0,
		Warnings:       []string{},
		RequiresReview: false,
	}

	// Track which expected fields have been mapped
	mappedFields := make(map[string]bool)

	// Process each header
	for i, header := range headers {
		// Get sample values for this column
		var sampleValues []string
		for _, row := range sampleData {
			if i < len(row) {
				sampleValues = append(sampleValues, strings.TrimSpace(row[i]))
			}
		}

		// Try to find a mapping for this header
		mapping, found := ds.fuzzyMatcher.FindBestMatch(header, ExpectedFields)

		if found {
			// Infer data type from sample data
			inferredType, _ := ds.dataTypeInferer.InferType(sampleValues, header)
			mapping.InferredType = inferredType
			mapping.SampleValues = ds.limitSampleValues(sampleValues, 3)

			// Check if field is already mapped
			if mappedFields[mapping.ExpectedField] {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Header '%s' maps to '%s' but this field is already mapped", header, mapping.ExpectedField))
				result.Unmapped = append(result.Unmapped, header)
				result.RequiresReview = true
				continue
			}

			// Validate mapping confidence
			if mapping.Confidence < 0.7 {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Low confidence mapping: '%s' -> '%s' (%.1f%% confidence)",
						header, mapping.ExpectedField, mapping.Confidence*100))
				result.RequiresReview = true
			}

			mappedFields[mapping.ExpectedField] = true
			result.Mappings = append(result.Mappings, *mapping)

		} else {
			result.Unmapped = append(result.Unmapped, header)
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Could not map header '%s' to any expected field", header))
			result.RequiresReview = true
		}
	}

	// Check for missing required fields
	result.checkRequiredFields()

	// Calculate overall confidence
	result.Confidence = ds.calculateOverallConfidence(result.Mappings)

	return result
}

// sampleRows reads a limited number of rows for sampling
func (ds *DetectionService) sampleRows(csvReader *csv.Reader, maxRows int) [][]string {
	var rows [][]string

	for i := 0; i < maxRows; i++ {
		row, err := csvReader.Read()
		if err != nil {
			break
		}
		rows = append(rows, row)
	}

	return rows
}

// limitSampleValues returns a limited number of sample values for display
func (ds *DetectionService) limitSampleValues(values []string, max int) []string {
	if len(values) <= max {
		return values
	}

	result := make([]string, max)
	copy(result, values[:max])

	// Add indicator if there are more values
	if len(values) > max {
		result = append(result, fmt.Sprintf("... and %d more", len(values)-max))
	}

	return result
}

// checkRequiredFields ensures all required fields are mapped
func (dr *DetectionResult) checkRequiredFields() {
	mappedFields := make(map[string]bool)
	for _, mapping := range dr.Mappings {
		mappedFields[mapping.ExpectedField] = true
	}

	for _, field := range ExpectedFields {
		if field.Required && !mappedFields[field.Name] {
			dr.Warnings = append(dr.Warnings,
				fmt.Sprintf("Required field '%s' is not mapped", field.Name))
			dr.RequiresReview = true
		}
	}
}

// calculateOverallConfidence calculates the overall confidence of all mappings
func (ds *DetectionService) calculateOverallConfidence(mappings []ColumnMapping) float64 {
	if len(mappings) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, mapping := range mappings {
		totalConfidence += mapping.Confidence
	}

	return totalConfidence / float64(len(mappings))
}

// ValidateMappings validates the detected mappings against expected fields
func (ds *DetectionService) ValidateMappings(mappings []ColumnMapping) []string {
	var errors []string

	// Check for duplicate field mappings
	fieldCount := make(map[string]int)
	for _, mapping := range mappings {
		fieldCount[mapping.ExpectedField]++
	}

	for field, count := range fieldCount {
		if count > 1 {
			errors = append(errors, fmt.Sprintf("Field '%s' is mapped multiple times", field))
		}
	}

	// Check for required fields
	mappedFields := make(map[string]bool)
	for _, mapping := range mappings {
		mappedFields[mapping.ExpectedField] = true
	}

	for _, field := range ExpectedFields {
		if field.Required && !mappedFields[field.Name] {
			errors = append(errors, fmt.Sprintf("Required field '%s' is not mapped", field.Name))
		}
	}

	return errors
}

// GetExpectedFields returns the list of expected fields for CSV import
func (ds *DetectionService) GetExpectedFields() []ExpectedField {
	return ExpectedFields
}