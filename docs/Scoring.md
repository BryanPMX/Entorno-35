# Scoring Implementation Guide – NOM-035 Risk Calculation

This document describes the scoring implementation for NOM-035-STPS-2018 psychosocial risk assessments: polarity rules (Guide II and III), risk level thresholds, and domain grouping. It is the reference for developers working on the scoring service or report logic.

**Last updated**: 2026-01-31

---

## Verified Data Source
This guide is based on rigorous verification against the official **NOM-035-STPS-2018** document (DOF - Diario Oficial de la Federación.pdf), specifically:
- Table 2 (Guide II polarity mapping)
- Table 3 (Guide II category/domain grouping)
- Table 5 (Guide III polarity mapping)
- Table 6 (Guide III category/domain grouping)
- Page 28 (Guide II thresholds)
- Page 36 (Guide III thresholds)

---

## 1. Polarity Scoring Rules

### Guide II (46 Questions)

#### POSITIVE Questions (0 → 4 scale)
**Question Numbers**: 18-33

**Scoring**:
- Siempre = 0 points
- Casi Siempre = 1 point
- Algunas Veces = 2 points
- Casi Nunca = 3 points
- Nunca = 4 points

**Interpretation**: Higher frequency of positive responses = Lower risk score

#### NEGATIVE Questions (4 → 0 scale)
**Question Numbers**: 1-17, 34-46

**Scoring**:
- Siempre = 4 points
- Casi Siempre = 3 points
- Algunas Veces = 2 points
- Casi Nunca = 1 point
- Nunca = 0 points

**Interpretation**: Higher frequency of negative responses = Higher risk score

---

### Guide III (72 Questions)

#### POSITIVE Questions (0 → 4 scale)
**Question Numbers**: 1, 4, 23-28, 30-41, 42-53, 55-57

**Scoring**: Same as Guide II POSITIVE
- Siempre = 0 points
- Casi Siempre = 1 point
- Algunas Veces = 2 points
- Casi Nunca = 3 points
- Nunca = 4 points

#### NEGATIVE Questions (4 → 0 scale)
**Question Numbers**: 2, 3, 5-22, 29, 54, 58-72

**Scoring**: Same as Guide II NEGATIVE
- Siempre = 4 points
- Casi Siempre = 3 points
- Algunas Veces = 2 points
- Casi Nunca = 1 point
- Nunca = 0 points

---

## 2. Risk Level Calculation Logic

### Critical Implementation Rule: "Less Than" Logic

The official NOM-035 standard uses **strict "less than" comparisons** for threshold boundaries.

#### Risk Level Determination Formula

```
If Score < Limit1:
    Risk Level = "Nulo"
Else If Limit1 <= Score < Limit2:
    Risk Level = "Bajo"
Else If Limit2 <= Score < Limit3:
    Risk Level = "Medio"
Else If Limit3 <= Score < Limit4:
    Risk Level = "Alto"
Else If Score >= Limit4:
    Risk Level = "Muy Alto"
```

#### Example: Guide III Total Score

Given ranges: `[50, 75, 99, 140]` (bajo, medio, alto, muy_alto thresholds)

```
If Score < 50:
    → Nulo
Else If 50 <= Score < 75:
    → Bajo (e.g., Score = 74.9 → Bajo)
Else If 75 <= Score < 99:
    → Medio (e.g., Score = 75.0 → Medio)
Else If 99 <= Score < 140:
    → Alto
Else If Score >= 140:
    → Muy Alto
```

**Key Point**: The boundary value belongs to the **higher** risk level.
- Score = 50.0 → Bajo (not Nulo)
- Score = 75.0 → Medio (not Bajo)
- Score = 99.0 → Alto (not Medio)
- Score = 140.0 → Muy Alto (not Alto)

---

## 3. Threshold Ranges Structure

### Guide II Thresholds (16-50 employees)

#### Total Score
```json
{
  "nulo": { "min": 0, "max": 20 },
  "bajo": { "min": 20, "max": 45 },
  "medio": { "min": 45, "max": 70 },
  "alto": { "min": 70, "max": 90 },
  "muy_alto": { "min": 90, "max": 999 }
}
```

**Ranges Array**: `[20, 45, 70, 90]`
- Index 0 (20) = Bajo threshold
- Index 1 (45) = Medio threshold
- Index 2 (70) = Alto threshold
- Index 3 (90) = Muy Alto threshold

#### Categories
1. **Ambiente de trabajo**: `[3, 5, 7, 9]`
2. **Factores propios de la actividad**: `[10, 20, 30, 40]`
3. **Organización del tiempo de trabajo**: `[4, 6, 9, 12]`
4. **Liderazgo y relaciones en el trabajo**: `[10, 18, 28, 38]`

#### Domains
1. **Condiciones en el ambiente de trabajo**: `[3, 5, 7, 9]`
2. **Carga de trabajo**: `[12, 16, 20, 24]`
3. **Falta de control sobre el trabajo**: `[5, 8, 11, 14]`
4. **Jornada de trabajo**: `[1, 2, 4, 6]`
5. **Interferencia en la relación trabajo-familia**: `[1, 2, 4, 6]`
6. **Liderazgo**: `[3, 5, 8, 11]`
7. **Relaciones en el trabajo**: `[5, 8, 11, 14]`
8. **Violencia**: `[7, 10, 13, 16]`

---

### Guide III Thresholds (>50 employees)

#### Total Score
```json
{
  "nulo": { "min": 0, "max": 50 },
  "bajo": { "min": 50, "max": 75 },
  "medio": { "min": 75, "max": 99 },
  "alto": { "min": 99, "max": 140 },
  "muy_alto": { "min": 140, "max": 999 }
}
```

**Ranges Array**: `[50, 75, 99, 140]`

#### Categories (5 total - includes "Entorno organizacional")
1. **Ambiente de trabajo**: `[5, 9, 11, 14]`
2. **Factores propios de la actividad**: `[15, 30, 45, 60]`
3. **Organización del tiempo de trabajo**: `[5, 7, 10, 13]`
4. **Liderazgo y relaciones en el trabajo**: `[14, 29, 42, 58]`
5. **Entorno organizacional**: `[10, 14, 18, 23]` ⭐ (Guide III only)

#### Domains (10 total - Guide III has 2 additional domains)
1. **Condiciones en el ambiente de trabajo**: `[5, 9, 11, 14]`
2. **Carga de trabajo**: `[15, 21, 27, 37]`
3. **Falta de control sobre el trabajo**: `[11, 16, 21, 25]`
4. **Jornada de trabajo**: `[1, 2, 4, 6]`
5. **Interferencia en la relación trabajo-familia**: `[4, 6, 8, 10]`
6. **Liderazgo**: `[9, 12, 16, 20]`
7. **Relaciones en el trabajo**: `[10, 13, 17, 21]`
8. **Violencia**: `[7, 10, 13, 16]`
9. **Reconocimiento del desempeño**: `[6, 10, 14, 18]` ⭐ (Guide III only)
10. **Insuficiente sentido de pertenencia e inestabilidad**: `[4, 6, 8, 10]` ⭐ (Guide III only)

---

## 4. Go Implementation Example

### GetRiskLabel Function

```go
package scoring

import (
    "math"
)

// RiskLevel represents the five risk levels
type RiskLevel string

const (
    RiskNulo    RiskLevel = "Nulo"
    RiskBajo    RiskLevel = "Bajo"
    RiskMedio   RiskLevel = "Medio"
    RiskAlto    RiskLevel = "Alto"
    RiskMuyAlto RiskLevel = "Muy Alto"
)

// GetRiskLabel calculates the risk level based on score and thresholds
// ranges: [bajoThreshold, medioThreshold, altoThreshold, muyAltoThreshold]
func GetRiskLabel(score float64, ranges []int) RiskLevel {
    // Nulo: score < ranges[0]
    if score < float64(ranges[0]) {
        return RiskNulo
    }
    
    // Bajo: ranges[0] <= score < ranges[1]
    if score >= float64(ranges[0]) && score < float64(ranges[1]) {
        return RiskBajo
    }
    
    // Medio: ranges[1] <= score < ranges[2]
    if score >= float64(ranges[1]) && score < float64(ranges[2]) {
        return RiskMedio
    }
    
    // Alto: ranges[2] <= score < ranges[3]
    if score >= float64(ranges[2]) && score < float64(ranges[3]) {
        return RiskAlto
    }
    
    // Muy Alto: score >= ranges[3]
    return RiskMuyAlto
}

// ApplyPolarity applies polarity inversion based on question polarity
// For POSITIVE: Siempre=0, Nunca=4 (inverted)
// For NEGATIVE: Siempre=4, Nunca=0 (direct)
func ApplyPolarity(responseValue int, polarity string) int {
    if polarity == "POSITIVE" {
        // Invert: 0->4, 1->3, 2->2, 3->1, 4->0
        return 4 - responseValue
    }
    // NEGATIVE: direct mapping
    return responseValue
}
```

### Test Cases

```go
func TestGetRiskLabel(t *testing.T) {
    guideIII := []int{50, 75, 99, 140}
    
    // Test boundary values
    assert.Equal(t, RiskNulo, GetRiskLabel(49.9, guideIII))
    assert.Equal(t, RiskBajo, GetRiskLabel(50.0, guideIII))  // Boundary → Bajo
    assert.Equal(t, RiskBajo, GetRiskLabel(74.9, guideIII))
    assert.Equal(t, RiskMedio, GetRiskLabel(75.0, guideIII)) // Boundary → Medio
    assert.Equal(t, RiskMedio, GetRiskLabel(98.9, guideIII))
    assert.Equal(t, RiskAlto, GetRiskLabel(99.0, guideIII))  // Boundary → Alto
    assert.Equal(t, RiskAlto, GetRiskLabel(139.9, guideIII))
    assert.Equal(t, RiskMuyAlto, GetRiskLabel(140.0, guideIII)) // Boundary → Muy Alto
    assert.Equal(t, RiskMuyAlto, GetRiskLabel(200.0, guideIII))
}
```

---

## 5. Critical Implementation Notes

### WARNING: Must Follow These Rules:

1. **Polarity MUST be applied at question level BEFORE aggregation**
   - Wrong: Sum responses, then invert
   - Correct: Invert each question score, then sum

2. **Use strict "less than" comparisons for boundaries**
   - Boundary values belong to the higher risk level
   - Score = threshold → Higher level (e.g., 75.0 → Medio, not Bajo)

3. **Guide II vs Guide III differences**:
   - Different total score thresholds
   - Guide III has 5 categories (vs 4 in Guide II)
   - Guide III has 10 domains (vs 8 in Guide II)
   - Guide III includes "Entorno organizacional" category

4. **Floating point precision**:
   - Use `float64` for calculations
   - Round to appropriate decimal places for display
   - Be aware of floating point comparison edge cases

5. **Data source**: Always use `RiskStrategy.json` for threshold values
   - Load at application startup
   - Cache in memory for performance
   - Version control the JSON file

---

## 6. Files Reference

- **Questions**: `nom035_questions.json` (138 questions with metadata)
- **Thresholds**: `RiskStrategy.json` (verified thresholds for both guides)
- **Source**: Official NOM-035-STPS-2018 PDF (DOF - Diario Oficial de la Federación)

---

**Last Verified**: Based on rigorous review of official NOM-035-STPS-2018 document  
**Status**: Ready for Phase 2 Implementation (Scoring Engine)

