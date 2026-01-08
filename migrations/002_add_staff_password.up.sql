-- Add password_hash column to staff table for staff authentication
ALTER TABLE staff ADD COLUMN password_hash VARCHAR(255);

-- Add comment for documentation
COMMENT ON COLUMN staff.password_hash IS 'Bcrypt hashed password for staff authentication';