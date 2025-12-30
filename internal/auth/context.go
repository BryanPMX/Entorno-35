package auth

import "github.com/google/uuid"

// Context represents the authenticated user's context (high cohesion - auth-specific data)
type Context struct {
	CompanyID uuid.UUID
	StaffID   *uuid.UUID // Optional, may be nil for company-level operations
	Email     string
	Role      string
}

// IsStaff returns true if the context has a staff ID
func (c *Context) IsStaff() bool {
	return c.StaffID != nil
}

// GetStaffID returns the staff ID or panics if not set
// Use IsStaff() to check before calling this
func (c *Context) GetStaffID() uuid.UUID {
	if c.StaffID == nil {
		panic("staff ID not set in auth context")
	}
	return *c.StaffID
}

