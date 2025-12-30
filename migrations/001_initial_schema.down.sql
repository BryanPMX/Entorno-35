-- Drop triggers
DROP TRIGGER IF EXISTS update_responses_updated_at ON responses;
DROP TRIGGER IF EXISTS update_assessments_updated_at ON assessments;
DROP TRIGGER IF EXISTS update_assessment_links_updated_at ON assessment_links;
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS update_dimensions_updated_at ON dimensions;
DROP TRIGGER IF EXISTS update_domains_updated_at ON domains;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_staff_updated_at ON staff;
DROP TRIGGER IF EXISTS update_companies_updated_at ON companies;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS responses;
DROP TABLE IF EXISTS assessment_links;
DROP TABLE IF EXISTS assessments;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS dimensions;
DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS staff;
DROP TABLE IF EXISTS companies;

-- Drop custom types/enums
DROP TYPE IF EXISTS risk_level;
DROP TYPE IF EXISTS assessment_status;
DROP TYPE IF EXISTS guide_type;
DROP TYPE IF EXISTS question_polarity;
DROP TYPE IF EXISTS question_type;
DROP TYPE IF EXISTS subscription_status;

