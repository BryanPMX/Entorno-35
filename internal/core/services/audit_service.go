package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AuditService handles file-based audit logging for critical operations
// Logs are written to a file that is excluded from git commits
type AuditService struct {
	logFile   *os.File
	logMutex  sync.Mutex
	logPath   string
}

var (
	auditServiceInstance *AuditService
	auditServiceOnce     sync.Once
)

// GetAuditService returns the singleton audit service instance
func GetAuditService() *AuditService {
	auditServiceOnce.Do(func() {
		// Create logs directory if it doesn't exist
		logsDir := "logs"
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			// If we can't create logs directory, fall back to stdout
			fmt.Printf("[AUDIT WARNING] Failed to create logs directory: %v. Using stdout.\n", err)
			auditServiceInstance = &AuditService{
				logPath: "",
			}
			return
		}

		// Create audit log file with date suffix
		logFileName := fmt.Sprintf("audit_%s.log", time.Now().Format("2006-01-02"))
		logPath := filepath.Join(logsDir, logFileName)

		// Open file in append mode, create if doesn't exist
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// If we can't open file, fall back to stdout
			fmt.Printf("[AUDIT WARNING] Failed to open audit log file: %v. Using stdout.\n", err)
			auditServiceInstance = &AuditService{
				logPath: "",
			}
			return
		}

		auditServiceInstance = &AuditService{
			logFile: file,
			logPath: logPath,
		}
	})

	return auditServiceInstance
}

// LogAction logs an audit action to the audit log file
func (s *AuditService) LogAction(action string, companyID, entityID uuid.UUID, details string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[AUDIT] Action=%s CompanyID=%s EntityID=%s Details=%s Time=%s\n",
		action, companyID.String(), entityID.String(), details, timestamp)

	// Try to write to file
	if s.logFile != nil {
		if _, err := s.logFile.WriteString(logEntry); err != nil {
			// If file write fails, fall back to stdout
			fmt.Printf("[AUDIT FALLBACK] %s", logEntry)
		} else {
			// Flush to ensure data is written
			s.logFile.Sync()
		}
	} else {
		// Fall back to stdout if file is not available
		fmt.Printf("[AUDIT] %s", logEntry)
	}
}

// LogAssessmentAction logs an assessment-related audit action
func (s *AuditService) LogAssessmentAction(action string, companyID, assessmentID, staffID uuid.UUID, details string) {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[AUDIT] Action=%s CompanyID=%s AssessmentID=%s StaffID=%s Details=%s Time=%s\n",
		action, companyID.String(), assessmentID.String(), staffID.String(), details, timestamp)

	// Try to write to file
	if s.logFile != nil {
		if _, err := s.logFile.WriteString(logEntry); err != nil {
			// If file write fails, fall back to stdout
			fmt.Printf("[AUDIT FALLBACK] %s", logEntry)
		} else {
			// Flush to ensure data is written
			s.logFile.Sync()
		}
	} else {
		// Fall back to stdout if file is not available
		fmt.Printf("[AUDIT] %s", logEntry)
	}
}

// Close closes the audit log file
func (s *AuditService) Close() error {
	s.logMutex.Lock()
	defer s.logMutex.Unlock()

	if s.logFile != nil {
		return s.logFile.Close()
	}
	return nil
}
