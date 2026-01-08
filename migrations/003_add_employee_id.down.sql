-- Revert employee_id addition and CURP nullability changes

-- Remove employee_id column
ALTER TABLE staff DROP COLUMN employee_id;

-- Make CURP NOT NULL again
ALTER TABLE staff ALTER COLUMN curp SET NOT NULL;