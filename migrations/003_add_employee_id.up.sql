-- Add employee_id field for staff without CURPs and make CURP nullable
-- This enables hybrid authentication: CURP + password OR employee_id + password

ALTER TABLE staff ADD COLUMN employee_id VARCHAR(50) UNIQUE;

-- Make CURP nullable to support staff without CURPs
ALTER TABLE staff ALTER COLUMN curp DROP NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN staff.employee_id IS 'Auto-generated unique identifier for staff without CURPs. Format: CMP{COMPANY_ID}-{SEQUENCE}';
COMMENT ON COLUMN staff.curp IS 'Clave Única de Registro de Población (optional for staff without CURPs)';