-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create custom types/enums
CREATE TYPE subscription_status AS ENUM ('active', 'inactive');
CREATE TYPE question_type AS ENUM ('binary', 'likert');
CREATE TYPE question_polarity AS ENUM ('POSITIVE', 'NEGATIVE');
CREATE TYPE guide_type AS ENUM ('I', 'II', 'III');
CREATE TYPE assessment_status AS ENUM ('pending', 'completed', 'cancelled');
CREATE TYPE risk_level AS ENUM ('nulo', 'bajo', 'medio', 'alto', 'muy_alto');

-- Companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rfc VARCHAR(13) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    subscription_status subscription_status NOT NULL DEFAULT 'inactive',
    employee_count INTEGER NOT NULL DEFAULT 0,
    subscription_start_date DATE,
    subscription_end_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_companies_rfc ON companies(rfc);
CREATE INDEX idx_companies_subscription_status ON companies(subscription_status);
CREATE INDEX idx_companies_deleted_at ON companies(deleted_at);

-- Staff table
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    curp VARCHAR(18) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    demographics JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Unique constraint: CURP must be unique per company
    UNIQUE(company_id, curp)
);

CREATE INDEX idx_staff_company_id ON staff(company_id);
CREATE INDEX idx_staff_curp ON staff(curp);
CREATE INDEX idx_staff_company_curp ON staff(company_id, curp);
CREATE INDEX idx_staff_deleted_at ON staff(deleted_at);
CREATE INDEX idx_staff_demographics ON staff USING GIN(demographics);

-- Categories table (normalized for reporting)
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_name ON categories(name);

-- Domains table (normalized for reporting)
CREATE TABLE domains (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, name)
);

CREATE INDEX idx_domains_category_id ON domains(category_id);

-- Dimensions table (normalized for reporting)
CREATE TABLE dimensions (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(domain_id, name)
);

CREATE INDEX idx_dimensions_domain_id ON dimensions(domain_id);

-- Questions table
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    question_number INTEGER NOT NULL,
    guide_type guide_type NOT NULL,
    type question_type NOT NULL,
    text TEXT NOT NULL,
    polarity question_polarity NOT NULL,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE RESTRICT,
    dimension_id INTEGER REFERENCES dimensions(id) ON DELETE RESTRICT,
    order_index INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_questions_guide_type ON questions(guide_type);
CREATE INDEX idx_questions_category_id ON questions(category_id);
CREATE INDEX idx_questions_domain_id ON questions(domain_id);
CREATE INDEX idx_questions_dimension_id ON questions(dimension_id);
CREATE INDEX idx_questions_guide_number ON questions(guide_type, question_number);

-- Assessments table (created before assessment_links due to FK reference)
CREATE TABLE assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    period INTEGER NOT NULL,
    guide_type guide_type NOT NULL,
    status assessment_status NOT NULL DEFAULT 'pending',
    total_score DECIMAL(10, 2),
    risk_level risk_level,
    requires_medical_attention BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_assessments_staff_id ON assessments(staff_id);
CREATE INDEX idx_assessments_company_id ON assessments(company_id);
CREATE INDEX idx_assessments_period ON assessments(period);
CREATE INDEX idx_assessments_status ON assessments(status);
CREATE INDEX idx_assessments_staff_period ON assessments(staff_id, period);
CREATE INDEX idx_assessments_deleted_at ON assessments(deleted_at);

-- Responses table
CREATE TABLE responses (
    id SERIAL PRIMARY KEY,
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    selected_value INTEGER NOT NULL CHECK (selected_value >= 0 AND selected_value <= 4),
    calculated_score INTEGER NOT NULL,
    answered_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Unique constraint: One response per question per assessment
    UNIQUE(assessment_id, question_id)
);

CREATE INDEX idx_responses_assessment_id ON responses(assessment_id);
CREATE INDEX idx_responses_question_id ON responses(question_id);
CREATE INDEX idx_responses_assessment_question ON responses(assessment_id, question_id);

-- Assessment Links table (created after assessments due to FK reference)
CREATE TABLE assessment_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) NOT NULL UNIQUE,
    staff_id UUID NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    assessment_id UUID REFERENCES assessments(id) ON DELETE SET NULL,
    expires_at TIMESTAMP NOT NULL,
    accessed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_assessment_links_token ON assessment_links(token);
CREATE INDEX idx_assessment_links_staff_id ON assessment_links(staff_id);
CREATE INDEX idx_assessment_links_assessment_id ON assessment_links(assessment_id);
CREATE INDEX idx_assessment_links_expires_at ON assessment_links(expires_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_staff_updated_at BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_domains_updated_at BEFORE UPDATE ON domains
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_dimensions_updated_at BEFORE UPDATE ON dimensions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_questions_updated_at BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assessment_links_updated_at BEFORE UPDATE ON assessment_links
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assessments_updated_at BEFORE UPDATE ON assessments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_responses_updated_at BEFORE UPDATE ON responses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

