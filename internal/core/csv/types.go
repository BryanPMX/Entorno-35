package csv

import (
	"regexp"
	"strings"
)

// ExpectedField represents a field that the system expects in CSV imports
type ExpectedField struct {
	Name        string         `json:"name"`
	Type        FieldType      `json:"type"`
	Required    bool           `json:"required"`
	Description string         `json:"description"`
	Validator   FieldValidator `json:"-"`
	Synonyms    []string       `json:"synonyms"`
}

// FieldType represents the data type of a field
type FieldType string

const (
	FieldTypeText   FieldType = "text"
	FieldTypeNumber FieldType = "number"
	FieldTypeEmail  FieldType = "email"
	FieldTypeEnum   FieldType = "enum"
)

// FieldValidator validates field values
type FieldValidator interface {
	Validate(value string) bool
}

// ColumnMapping represents a detected mapping between CSV header and expected field
type ColumnMapping struct {
	CSVHeader      string        `json:"csv_header"`
	ExpectedField  string        `json:"expected_field"`
	Confidence     float64       `json:"confidence"`
	MatchType      MatchType     `json:"match_type"`
	SampleValues   []string      `json:"sample_values,omitempty"`
	InferredType   FieldType     `json:"inferred_type,omitempty"`
}

// MatchType indicates how the mapping was determined
type MatchType string

const (
	MatchTypeExact     MatchType = "exact"      // Perfect match
	MatchTypeSynonym   MatchType = "synonym"    // Matched via synonym dictionary
	MatchTypeFuzzy     MatchType = "fuzzy"      // Fuzzy string matching
	MatchTypeInference MatchType = "inference"  // Based on data analysis
	MatchTypeManual    MatchType = "manual"     // User-specified mapping
)

// DetectionResult contains the complete analysis of a CSV file
type DetectionResult struct {
	Mappings       []ColumnMapping `json:"mappings"`
	Unmapped       []string        `json:"unmapped_headers"`
	Confidence     float64         `json:"overall_confidence"`
	Warnings       []string        `json:"warnings"`
	RequiresReview bool            `json:"requires_review"`
}

// Standard expected fields for staff import
var ExpectedFields = []ExpectedField{
	{
		Name:        "name",
		Type:        FieldTypeText,
		Required:    true,
		Description: "Staff member's full name",
		Validator:   &TextValidator{MinLength: 2, MaxLength: 255},
		Synonyms:    []string{"full_name", "full name", "employee_name", "employee name", "staff_name", "staff name", "nombre", "nombre_completo", "nombre completo"},
	},
	{
		Name:        "curp",
		Type:        FieldTypeText,
		Required:    false,
		Description: "Clave Única de Registro de Población",
		Validator:   &CURPValidator{},
		Synonyms:    []string{"rfc", "id", "employee_id", "employee id", "personal_id", "identificacion"},
	},
	{
		Name:        "email",
		Type:        FieldTypeEmail,
		Required:    false,
		Description: "Email address",
		Validator:   &EmailValidator{},
		Synonyms:    []string{"email_address", "email address", "correo", "correo_electronico", "correo electronico", "mail"},
	},
	{
		Name:        "area",
		Type:        FieldTypeText,
		Required:    false,
		Description: "Department or area",
		Validator:   &TextValidator{MinLength: 1, MaxLength: 100},
		Synonyms:    []string{"department", "departamento", "area", "sector", "division"},
	},
	{
		Name:        "job",
		Type:        FieldTypeText,
		Required:    false,
		Description: "Job title or position",
		Validator:   &TextValidator{MinLength: 1, MaxLength: 100},
		Synonyms:    []string{"position", "role", "cargo", "puesto", "titulo", "title", "ocupacion"},
	},
	{
		Name:        "shift",
		Type:        FieldTypeEnum,
		Required:    false,
		Description: "Work shift",
		Validator:   &EnumValidator{AllowedValues: []string{"diurno", "nocturno", "mixto", "day", "night", "mixed", "regular"}},
		Synonyms:    []string{"turno", "jornada", "schedule", "horario"},
	},
	{
		Name:        "gender",
		Type:        FieldTypeEnum,
		Required:    false,
		Description: "Gender",
		Validator:   &EnumValidator{AllowedValues: []string{"masculino", "femenino", "otro", "male", "female", "other", "m", "f"}},
		Synonyms:    []string{"sexo", "genero", "sex"},
	},
}

// Validators

type TextValidator struct {
	MinLength int
	MaxLength int
}

func (v *TextValidator) Validate(value string) bool {
	length := len(value)
	return length >= v.MinLength && length <= v.MaxLength
}

type CURPValidator struct{}

func (v *CURPValidator) Validate(value string) bool {
	// CURP must be exactly 18 characters when provided
	return len(value) == 0 || len(value) == 18
}

type EmailValidator struct{}

func (v *EmailValidator) Validate(value string) bool {
	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}

type EnumValidator struct {
	AllowedValues []string
}

func (v *EnumValidator) Validate(value string) bool {
	if value == "" {
		return true // Empty values are allowed for optional fields
	}
	normalizedValue := strings.ToLower(strings.TrimSpace(value))
	for _, allowed := range v.AllowedValues {
		if normalizedValue == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}