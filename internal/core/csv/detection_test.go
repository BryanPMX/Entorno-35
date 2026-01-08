package csv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectionService_AnalyzeCSV(t *testing.T) {
	detector := NewDetectionService()

	// Test CSV with Spanish headers
	csvData := `Nombre Completo,CURP,Correo Electrónico,Departamento,Puesto,Turno,Género
Juan Pérez García,ABCD123456HIJKLM01,juan@example.com,Producción,Operador,Diurno,Masculino
María González López,EFGH567890MNOPQR02,maria@example.com,Recursos Humanos,Analista,Diurno,Femenino
Carlos Rodríguez,JKLM901234STUVWX03,carlos@example.com,Calidad,Supervisor,Nocturno,Masculino`

	reader := strings.NewReader(csvData)

	result, err := detector.AnalyzeCSV(reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Check that mappings were created
	assert.Greater(t, len(result.Mappings), 0)

	// Check that name field was mapped
	foundName := false
	for _, mapping := range result.Mappings {
		if mapping.ExpectedField == "name" {
			foundName = true
			assert.Greater(t, mapping.Confidence, 0.8) // Should have high confidence
			assert.Equal(t, MatchTypeSynonym, mapping.MatchType)
		}
	}
	assert.True(t, foundName, "Name field should be mapped")

	// Check that overall confidence is reasonable
	assert.Greater(t, result.Confidence, 0.5)
}

func TestFuzzyMatcher_FindBestMatch(t *testing.T) {
	matcher := NewFuzzyMatcher()

	tests := []struct {
		header   string
		expected string
		minConf  float64
	}{
		{"Full Name", "name", 0.9},
		{"Nombre Completo", "name", 0.9},
		{"Correo", "email", 0.8},
		{"Email Address", "email", 0.8},
		{"Department", "area", 0.8},
		{"Área", "area", 0.8},
		{"Puesto", "job", 0.8},
		{"Position", "job", 0.8},
		{"Turno", "shift", 0.8},
		{"Shift", "shift", 0.8},
		{"Género", "gender", 0.8},
		{"Gender", "gender", 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			mapping, found := matcher.FindBestMatch(tt.header, ExpectedFields)
			assert.True(t, found, "Should find mapping for %s", tt.header)
			assert.Equal(t, tt.expected, mapping.ExpectedField)
			assert.GreaterOrEqual(t, mapping.Confidence, tt.minConf)
		})
	}
}

func TestDataTypeInferer_InferType(t *testing.T) {
	inferer := NewDataTypeInferer()

	tests := []struct {
		samples []string
		header  string
		expected FieldType
	}{
		{[]string{"juan@example.com", "maria@test.com"}, "email", FieldTypeEmail},
		{[]string{"ABCD123456HIJKLM01", "EFGH567890MNOPQR02"}, "curp", FieldTypeNumber},
		{[]string{"Juan Pérez", "María García"}, "name", FieldTypeText},
		{[]string{"Masculino", "Femenino", "Masculino"}, "gender", FieldTypeEnum},
		{[]string{"Diurno", "Nocturno", "Mixto"}, "shift", FieldTypeEnum},
	}

	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			inferred, confidence := inferer.InferType(tt.samples, tt.header)
			assert.Equal(t, tt.expected, inferred)
			assert.Greater(t, confidence, 0.5)
		})
	}
}