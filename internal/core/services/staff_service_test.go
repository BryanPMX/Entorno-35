package services

import (
	"fmt"
	"strings"
	"testing"

	"github.com/entorno35/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockStaffRepository is a mock implementation of StaffRepository for testing
type MockStaffRepository struct {
	mock.Mock
}

func (m *MockStaffRepository) GetByIDAndCompany(id uuid.UUID, companyID uuid.UUID) (*domain.Staff, error) {
	args := m.Called(id, companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Staff), args.Error(1)
}

func (m *MockStaffRepository) Create(staff *domain.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *MockStaffRepository) Update(staff *domain.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *MockStaffRepository) Delete(id uuid.UUID, companyID uuid.UUID) error {
	args := m.Called(id, companyID)
	return args.Error(0)
}

func (m *MockStaffRepository) ListByCompany(companyID uuid.UUID, limit, offset int) ([]domain.Staff, int64, error) {
	args := m.Called(companyID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Staff), args.Get(1).(int64), args.Error(2)
}

func (m *MockStaffRepository) BulkCreate(staff []*domain.Staff) (int, error) {
	args := m.Called(staff)
	return args.Int(0), args.Error(1)
}

func TestStaffService_ImportFromCSV_ValidCSV(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez García,ABCD123456HIJKLM01,juan.perez@example.com,Producción,Operador,Diurno,Masculino
María González López,EFGH567890MNOPQR02,maria.gonzalez@example.com,Recursos Humanos,Analista,Diurno,Femenino
`

	mockRepo.On("BulkCreate", mock.Anything).Return(2, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalProcessed)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.SkippedCount)
	assert.Empty(t, result.Errors)

	// Verify BulkCreate was called with correct data
	call := mockRepo.Calls[0]
	staffList := call.Arguments[0].([]*domain.Staff)

	require.Len(t, staffList, 2)
	assert.Equal(t, "ABCD123456HIJKLM01", staffList[0].CURP)
	assert.Equal(t, "Juan Pérez García", staffList[0].FullName)
	assert.Equal(t, "juan.perez@example.com", staffList[0].Email)
	assert.Equal(t, "Producción", staffList[0].Demographics.Department)
	assert.Equal(t, "Operador", staffList[0].Demographics.Role)
	assert.Equal(t, "Diurno", staffList[0].Demographics.ShiftType)
	assert.Equal(t, "Masculino", staffList[0].Demographics.Gender)

	mockRepo.AssertExpectations(t)
}

func TestStaffService_ImportFromCSV_MissingHeader(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area
Juan Pérez,ABCD123456HIJKLM01,juan@example.com,Producción
`

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required header")
	assert.Nil(t, result)
	mockRepo.AssertNotCalled(t, "BulkCreate")
}

func TestStaffService_ImportFromCSV_InvalidCURP_TooShort(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,SHORT,email@example.com,Area,Job,Shift,Gender
`

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalProcessed)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.SkippedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "CURP must be exactly 18 characters")

	mockRepo.On("BulkCreate", mock.Anything).Return(0, nil)
	mockRepo.AssertNotCalled(t, "BulkCreate")
}

func TestStaffService_ImportFromCSV_InvalidCURP_TooLong(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,ABCD123456HIJKLM012345,email@example.com,Area,Job,Shift,Gender
`

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalProcessed)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.SkippedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "CURP must be exactly 18 characters")
}

func TestStaffService_ImportFromCSV_MissingName(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
,ABCD123456HIJKLM01,email@example.com,Area,Job,Shift,Gender
`

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalProcessed)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.SkippedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "Name is required")
}

func TestStaffService_ImportFromCSV_MissingCURP(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,,email@example.com,Area,Job,Shift,Gender
`

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalProcessed)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.SkippedCount)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "CURP is required")
}

func TestStaffService_ImportFromCSV_CURP_Normalization(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	// CURP with lowercase and spaces
	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,  abcd123456hijklm01  ,email@example.com,Area,Job,Shift,Gender
`

	mockRepo.On("BulkCreate", mock.MatchedBy(func(staff []*domain.Staff) bool {
		return len(staff) == 1 && staff[0].CURP == "ABCD123456HIJKLM01"
	})).Return(1, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	mockRepo.AssertExpectations(t)
}

func TestStaffService_ImportFromCSV_SpecialCharacters(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	// Names with special characters (Spanish accents, ñ)
	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
José María O'Brien,ABCD123456HIJKLM01,jose@example.com,Producción,Operador,Diurno,Masculino
María José Ñoño,EFGH567890MNOPQR02,maria@example.com,Recursos Humanos,Analista,Nocturno,Femenino
`

	mockRepo.On("BulkCreate", mock.Anything).Return(2, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Empty(t, result.Errors)

	call := mockRepo.Calls[0]
	staffList := call.Arguments[0].([]*domain.Staff)
	assert.Equal(t, "José María O'Brien", staffList[0].FullName)
	assert.Equal(t, "María José Ñoño", staffList[1].FullName)
}

func TestStaffService_ImportFromCSV_MixedValidInvalid(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Valid User 1,ABCD123456HIJKLM01,valid1@example.com,Area,Job,Shift,Gender
Invalid CURP,SHORT,invalid@example.com,Area,Job,Shift,Gender
Valid User 2,EFGH567890MNOPQR02,valid2@example.com,Area,Job,Shift,Gender
Missing Name,,missing@example.com,Area,Job,Shift,Gender
Valid User 3,IJKL901234RSTUVW03,valid3@example.com,Area,Job,Shift,Gender
`

	mockRepo.On("BulkCreate", mock.MatchedBy(func(staff []*domain.Staff) bool {
		return len(staff) == 3 // Only valid records
	})).Return(3, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 5, result.TotalProcessed)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 2, result.SkippedCount)
	assert.Len(t, result.Errors, 2)
	assert.Contains(t, result.Errors[0], "CURP must be exactly 18 characters")
	assert.Contains(t, result.Errors[1], "CURP is required")
}

func TestStaffService_ImportFromCSV_EmptyEmail(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,ABCD123456HIJKLM01,,Producción,Operador,Diurno,Masculino
`

	mockRepo.On("BulkCreate", mock.MatchedBy(func(staff []*domain.Staff) bool {
		return len(staff) == 1 && staff[0].Email == ""
	})).Return(1, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Empty(t, result.Errors)
}

func TestStaffService_ImportFromCSV_CaseInsensitiveHeaders(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	// Headers in different case
	csvData := `NAME,curp,EMAIL,area,Job,SHIFT,gender
Juan Pérez,ABCD123456HIJKLM01,juan@example.com,Producción,Operador,Diurno,Masculino
`

	mockRepo.On("BulkCreate", mock.Anything).Return(1, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Empty(t, result.Errors)
}

func TestStaffService_ImportFromCSV_WhitespaceInHeaders(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	// Headers with leading/trailing whitespace
	csvData := ` Name , CURP , Email , Area , Job , Shift , Gender 
Juan Pérez,ABCD123456HIJKLM01,juan@example.com,Producción,Operador,Diurno,Masculino
`

	mockRepo.On("BulkCreate", mock.Anything).Return(1, nil)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	require.NoError(t, err)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Empty(t, result.Errors)
}

func TestStaffService_ImportFromCSV_LargeDataset(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	// Generate CSV with 2500 rows (to test batch processing with multiple batches)
	var b strings.Builder
	b.WriteString("Name,CURP,Email,Area,Job,Shift,Gender\n")
	for i := 0; i < 2500; i++ {
		curp := generateCURP(i)
		// Use first 8 chars for email username (CURP is 18 chars, so [:8] is safe)
		b.WriteString(fmt.Sprintf("User %d,%s,user%d@example.com,Area,Job,Shift,Gender\n", i, curp, i))
	}

	// With batch processing, BulkCreate will be called multiple times (every 1000 rows)
	// First batch: 1000 rows, Second batch: 1000 rows, Third batch: 500 rows
	// Use separate On calls for each batch size
	mockRepo.On("BulkCreate", mock.MatchedBy(func(staff []*domain.Staff) bool {
		return len(staff) == 1000
	})).Return(1000, nil).Times(2) // First two batches

	mockRepo.On("BulkCreate", mock.MatchedBy(func(staff []*domain.Staff) bool {
		return len(staff) == 500
	})).Return(500, nil).Times(1) // Last batch

	result, err := service.ImportFromCSV(strings.NewReader(b.String()), companyID)

	require.NoError(t, err)
	assert.Equal(t, 2500, result.TotalProcessed)
	assert.Equal(t, 2500, result.SuccessCount)
	assert.Equal(t, 0, result.SkippedCount)
	assert.Empty(t, result.Errors)

	// Verify BulkCreate was called 3 times (2 batches of 1000 + 1 batch of 500)
	mockRepo.AssertNumberOfCalls(t, "BulkCreate", 3)
}

func TestStaffService_ImportFromCSV_BulkCreateFailure(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)
	companyID := uuid.New()

	csvData := `Name,CURP,Email,Area,Job,Shift,Gender
Juan Pérez,ABCD123456HIJKLM01,juan@example.com,Producción,Operador,Diurno,Masculino
`

	mockRepo.On("BulkCreate", mock.Anything).Return(0, assert.AnError)

	result, err := service.ImportFromCSV(strings.NewReader(csvData), companyID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Helper function to generate a valid-looking 18-character CURP for testing
// Format: 4 letters + 6 digits + 2 letters + 6 alphanumeric = 18 characters total
// Pattern: ^[A-Z]{4}[0-9]{6}[A-Z]{2}[A-Z0-9]{6}$
func generateCURP(index int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	
	curp := ""
	// First 4: letters
	for i := 0; i < 4; i++ {
		curp += string(letters[(index+i)%26])
	}
	// Next 6: digits
	for i := 0; i < 6; i++ {
		curp += string(digits[(index+i)%10])
	}
	// Next 2: letters
	curp += string(letters[index%26])
	curp += string(letters[(index+1)%26])
	// Next 6: alphanumeric (mix of letters and digits)
	for i := 0; i < 6; i++ {
		if (index+i)%2 == 0 {
			curp += string(digits[(index+i)%10])
		} else {
			curp += string(letters[(index+i)%26])
		}
	}
	
	return curp
}

