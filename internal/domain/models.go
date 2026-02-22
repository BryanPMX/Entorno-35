package domain

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================================================
// Enums
// ============================================================================

// SubscriptionStatus represents the subscription status of a company
type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
)

func (s SubscriptionStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SubscriptionStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan SubscriptionStatus")
	}
	*s = SubscriptionStatus(str)
	return nil
}

// QuestionType represents the type of question (Binary or Likert)
type QuestionType string

const (
	QuestionTypeBinary QuestionType = "binary"
	QuestionTypeLikert QuestionType = "likert"
)

func (q QuestionType) Value() (driver.Value, error) {
	return string(q), nil
}

func (q *QuestionType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan QuestionType")
	}
	*q = QuestionType(str)
	return nil
}

// QuestionPolarity represents the polarity for scoring (POSITIVE or NEGATIVE)
type QuestionPolarity string

const (
	QuestionPolarityPositive QuestionPolarity = "POSITIVE"
	QuestionPolarityNegative QuestionPolarity = "NEGATIVE"
)

func (p QuestionPolarity) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *QuestionPolarity) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan QuestionPolarity")
	}
	*p = QuestionPolarity(str)
	return nil
}

// GuideType represents which guide the question belongs to
type GuideType string

const (
	GuideTypeI   GuideType = "I"
	GuideTypeII  GuideType = "II"
	GuideTypeIII GuideType = "III"
)

func (g GuideType) Value() (driver.Value, error) {
	return string(g), nil
}

func (g *GuideType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan GuideType")
	}
	*g = GuideType(str)
	return nil
}

// AssessmentStatus represents the status of an assessment
type AssessmentStatus string

const (
	AssessmentStatusPending   AssessmentStatus = "pending"
	AssessmentStatusCompleted AssessmentStatus = "completed"
	AssessmentStatusCancelled AssessmentStatus = "cancelled"
)

func (a AssessmentStatus) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *AssessmentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan AssessmentStatus")
	}
	*a = AssessmentStatus(str)
	return nil
}

// RiskLevel represents the calculated risk level
type RiskLevel string

const (
	RiskLevelNulo    RiskLevel = "nulo"
	RiskLevelBajo    RiskLevel = "bajo"
	RiskLevelMedio   RiskLevel = "medio"
	RiskLevelAlto    RiskLevel = "alto"
	RiskLevelMuyAlto RiskLevel = "muy_alto"
)

func (r RiskLevel) Value() (driver.Value, error) {
	return string(r), nil
}

func (r *RiskLevel) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan RiskLevel")
	}
	*r = RiskLevel(str)
	return nil
}

// ============================================================================
// Models
// ============================================================================

// Company represents a tenant company in the multi-tenant system
type Company struct {
	ID                    uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RFC                   string             `gorm:"type:varchar(13);uniqueIndex;not null" json:"rfc"`
	Name                  string             `gorm:"type:varchar(255);not null" json:"name"`
	Address               string             `gorm:"type:text" json:"address"`
	AdminEmail            *string            `gorm:"type:varchar(255);index" json:"admin_email,omitempty"`
	AdminPasswordHash     *string            `gorm:"type:varchar(255)" json:"-"`
	SubscriptionStatus    SubscriptionStatus `gorm:"type:varchar(20);not null;default:'inactive'" json:"subscription_status"`
	EmployeeCount         int                `gorm:"not null;default:0" json:"employee_count"`
	SubscriptionStartDate *time.Time         `gorm:"type:date" json:"subscription_start_date,omitempty"`
	SubscriptionEndDate   *time.Time         `gorm:"type:date" json:"subscription_end_date,omitempty"`
	StripeCustomerID      *string            `gorm:"type:varchar(255);index" json:"-"`
	StripeSubscriptionID  *string            `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	StripePriceID         *string            `gorm:"type:varchar(255)" json:"-"`
	CreatedAt             time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Staff       []Staff      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"staff,omitempty"`
	Assessments []Assessment `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"assessments,omitempty"`
}

// TableName specifies the table name for Company
func (Company) TableName() string {
	return "companies"
}

// BeforeCreate hook to generate UUID if not set
func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Staff represents a staff member (user) belonging to a company
type Staff struct {
	ID           uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CompanyID    uuid.UUID         `gorm:"type:uuid;not null;index" json:"company_id"`
	CURP         *string           `gorm:"type:varchar(18)" json:"curp,omitempty"`               // Made nullable for staff without CURPs
	EmployeeID   sql.NullString    `gorm:"type:varchar(50);unique" json:"employee_id,omitempty"` // Auto-generated for staff without CURPs, nullable
	FullName     string            `gorm:"type:varchar(255);not null" json:"full_name"`
	Email        string            `gorm:"type:varchar(255)" json:"email,omitempty"`
	PasswordHash string            `gorm:"type:varchar(255)" json:"-"` // Not exposed in JSON
	Demographics DemographicsJSONB `gorm:"type:jsonb" json:"demographics"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Company         Company          `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Assessments     []Assessment     `gorm:"foreignKey:StaffID;constraint:OnDelete:CASCADE" json:"assessments,omitempty"`
	AssessmentLinks []AssessmentLink `gorm:"foreignKey:StaffID;constraint:OnDelete:CASCADE" json:"assessment_links,omitempty"`
}

// TableName specifies the table name for Staff
func (Staff) TableName() string {
	return "staff"
}

// BeforeCreate hook to generate UUID if not set
func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// IsCURPAvailable returns true if the staff member has a CURP
func (s *Staff) IsCURPAvailable() bool {
	return s.CURP != nil && *s.CURP != ""
}

// GetIdentifier returns the primary identifier (CURP if available, otherwise employee_id)
func (s *Staff) GetIdentifier() string {
	if s.IsCURPAvailable() {
		return *s.CURP
	}
	if s.EmployeeID.Valid {
		return s.EmployeeID.String
	}
	return ""
}

// DemographicsJSONB represents the flexible demographics data structure
type DemographicsJSONB struct {
	Gender              string `json:"gender,omitempty"`                // e.g., "masculino", "femenino", "otro"
	AgeRange            string `json:"age_range,omitempty"`             // e.g., "18-25", "26-35", "36-45", "46-55", "56+"
	MaritalStatus       string `json:"marital_status,omitempty"`        // e.g., "soltero", "casado", "divorciado", "viudo"
	EducationLevel      string `json:"education_level,omitempty"`       // e.g., "secundaria", "preparatoria", "universidad", "postgrado"
	TimeInPosition      string `json:"time_in_position,omitempty"`      // e.g., "<1 año", "1-3 años", "3-5 años", "5+ años"
	ShiftType           string `json:"shift_type,omitempty"`            // e.g., "diurno", "nocturno", "mixto"
	ShiftRotation       string `json:"shift_rotation,omitempty"`        // e.g., "fijo", "rotativo"
	TotalWorkExperience string `json:"total_work_experience,omitempty"` // e.g., "<1 año", "1-5 años", "5-10 años", "10+ años"
	Department          string `json:"department,omitempty"`
	Role                string `json:"role,omitempty"`
}

// Value implements the driver.Valuer interface
func (d DemographicsJSONB) Value() (driver.Value, error) {
	return json.Marshal(d)
}

// Scan implements the sql.Scanner interface
func (d *DemographicsJSONB) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan DemographicsJSONB")
	}
	return json.Unmarshal(bytes, d)
}

// Category represents a question category (normalized table for reporting)
type Category struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Domains []Domain `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"domains,omitempty"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// Domain represents a question domain (normalized table for reporting)
type Domain struct {
	ID         uint      `gorm:"primary_key;autoIncrement" json:"id"`
	CategoryID uint      `gorm:"not null;index" json:"category_id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Category  Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Questions []Question `gorm:"foreignKey:DomainID;constraint:OnDelete:RESTRICT" json:"questions,omitempty"`
}

// TableName specifies the table name for Domain
func (Domain) TableName() string {
	return "domains"
}

// Dimension represents a question dimension (normalized table for reporting)
type Dimension struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	DomainID  uint      `gorm:"not null;index" json:"domain_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Domain    Domain     `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Questions []Question `gorm:"foreignKey:DimensionID;constraint:OnDelete:RESTRICT" json:"questions,omitempty"`
}

// TableName specifies the table name for Dimension
func (Dimension) TableName() string {
	return "dimensions"
}

// Question represents a question in the NOM-035 questionnaire
type Question struct {
	ID             uint         `gorm:"primary_key;autoIncrement" json:"id"`
	QuestionNumber int          `gorm:"not null" json:"question_number"`
	GuideType      GuideType    `gorm:"type:varchar(10);not null;index" json:"guide_type"`
	Type           QuestionType `gorm:"type:varchar(20);not null" json:"type"`
	Text           string       `gorm:"type:text;not null" json:"text"`
	// Polarity is only for Guide II/III (Likert scale questions)
	Polarity *QuestionPolarity `gorm:"type:varchar(10)" json:"polarity,omitempty"`
	// Section and Subsection are only for Guide I (Trauma assessment)
	Section    *string `gorm:"type:varchar(10)" json:"section,omitempty"`
	Subsection *string `gorm:"type:varchar(100)" json:"subsection,omitempty"`
	// Category, Domain, Dimension are only for Guide II/III
	CategoryID  *uint     `gorm:"index" json:"category_id,omitempty"`
	DomainID    *uint     `gorm:"index" json:"domain_id,omitempty"`
	DimensionID *uint     `gorm:"index" json:"dimension_id,omitempty"`
	OrderIndex  int       `gorm:"not null;default:0" json:"order_index"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships (optional for Guide I)
	Category  *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Domain    *Domain    `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Dimension *Dimension `gorm:"foreignKey:DimensionID" json:"dimension,omitempty"`
	Responses []Response `gorm:"foreignKey:QuestionID;constraint:OnDelete:RESTRICT" json:"responses,omitempty"`
}

// TableName specifies the table name for Question
func (Question) TableName() string {
	return "questions"
}

// AssessmentLink represents a secure link for staff to take assessments
type AssessmentLink struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Token        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	StaffID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"staff_id"`
	AssessmentID *uuid.UUID `gorm:"type:uuid;index" json:"assessment_id,omitempty"`
	ExpiresAt    time.Time  `gorm:"not null" json:"expires_at"`
	AccessedAt   *time.Time `gorm:"type:timestamp" json:"accessed_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Staff      Staff       `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	Assessment *Assessment `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
}

// TableName specifies the table name for AssessmentLink
func (AssessmentLink) TableName() string {
	return "assessment_links"
}

// BeforeCreate hook to generate UUID and token if not set
func (a *AssessmentLink) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Token == "" {
		a.Token = uuid.New().String()
	}
	return nil
}

// Assessment represents a completed or in-progress assessment
type Assessment struct {
	ID                       uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	StaffID                  uuid.UUID        `gorm:"type:uuid;not null;index" json:"staff_id"`
	CompanyID                uuid.UUID        `gorm:"type:uuid;not null;index" json:"company_id"`
	Period                   int              `gorm:"not null" json:"period"` // e.g., 2025
	GuideType                GuideType        `gorm:"type:varchar(10);not null" json:"guide_type"`
	Status                   AssessmentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	TotalScore               *float64         `gorm:"type:decimal(10,2)" json:"total_score,omitempty"`
	RiskLevel                *RiskLevel       `gorm:"type:varchar(20)" json:"risk_level,omitempty"`
	RequiresMedicalAttention bool             `gorm:"default:false" json:"requires_medical_attention"`
	CompletedAt              *time.Time       `gorm:"type:timestamp" json:"completed_at,omitempty"`
	CreatedAt                time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Staff     Staff      `gorm:"foreignKey:StaffID" json:"staff,omitempty"`
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Responses []Response `gorm:"foreignKey:AssessmentID;constraint:OnDelete:CASCADE" json:"responses,omitempty"`
}

// TableName specifies the table name for Assessment
func (Assessment) TableName() string {
	return "assessments"
}

// BeforeCreate hook to generate UUID if not set
func (a *Assessment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// Response represents a single response to a question in an assessment
type Response struct {
	ID              uint       `gorm:"primary_key;autoIncrement" json:"id"`
	AssessmentID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"assessment_id"`
	QuestionID      uint       `gorm:"not null;index" json:"question_id"`
	SelectedValue   int        `gorm:"not null;check:selected_value >= 0 AND selected_value <= 4" json:"selected_value"` // 0-4 for Likert, 0-1 for Binary
	CalculatedScore int        `gorm:"not null" json:"calculated_score"`                                                 // Score after polarity inversion
	AnsweredAt      *time.Time `gorm:"type:timestamp" json:"answered_at,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Assessment Assessment `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Question   Question   `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

// TableName specifies the table name for Response
func (Response) TableName() string {
	return "responses"
}

// Unique constraint: One response per question per assessment
// This will be added via migration
