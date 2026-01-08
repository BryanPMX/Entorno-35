# CSV Detection Package

Intelligent CSV column detection and mapping for staff import operations.

## Overview

This package provides intelligent auto-detection of CSV column mappings using:
- Fuzzy string matching algorithms
- Synonym recognition for common field variations
- Data type inference from sample data
- Confidence scoring for mapping suggestions

## Components

### Detection Service
- `DetectionService`: Main service for CSV analysis and mapping
- `ColumnMapper`: Maps detected columns to expected fields
- `DataTypeInferer`: Analyzes sample data to determine column types

### Fuzzy Matching
- Levenshtein distance algorithm
- Jaccard similarity for set-based matching
- Soundex phonetic matching for similar-sounding words

### Synonyms Dictionary
- Common field name variations (Name ↔ Full Name ↔ Employee Name)
- Multilingual support (Spanish/English)
- Domain-specific terminology

## Usage

```go
detector := csv.NewDetectionService()

// Analyze CSV headers
headers := []string{"Full Name", "Correo", "Departamento", "Puesto"}
mappings := detector.DetectColumnMappings(headers)

// Get confidence scores
for _, mapping := range mappings {
    fmt.Printf("%s -> %s (confidence: %.2f)\n",
        mapping.CSVHeader, mapping.ExpectedField, mapping.Confidence)
}
```

## Expected Fields

| Field | Type | Required | Synonyms |
|-------|------|----------|----------|
| Name | string | Yes | Full Name, Employee Name, Staff Name |
| CURP | string | No | RFC, ID, Employee ID |
| Email | string | No | Correo, Email Address |
| Area | string | No | Department, Área, Departamento |
| Job | string | No | Position, Role, Cargo, Puesto |
| Shift | string | No | Turno, Schedule, Jornada |
| Gender | string | No | Sexo, Género |

## Data Type Inference

Automatically detects:
- **Text**: String fields (names, emails)
- **Numbers**: CURP (18 chars), employee IDs
- **Enums**: Gender, Shift (limited values)
- **Emails**: RFC 5322 compliant patterns