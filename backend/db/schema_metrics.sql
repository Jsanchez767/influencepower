-- Official Performance Metrics Schema

-- Official metrics table (comprehensive performance data)
CREATE TABLE IF NOT EXISTS official_metrics (
    id BIGSERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) UNIQUE NOT NULL,
    
    -- Constituent Service Metrics
    avg_response_time_hours DECIMAL(10,2), -- Average hours to respond
    cases_resolved_total INTEGER DEFAULT 0,
    cases_resolved_last_month INTEGER DEFAULT 0,
    office_hours_held_total INTEGER DEFAULT 0,
    office_hours_held_last_month INTEGER DEFAULT 0,
    town_halls_attended_total INTEGER DEFAULT 0,
    town_halls_attended_last_year INTEGER DEFAULT 0,
    
    -- Legislative Productivity
    bills_introduced_total INTEGER DEFAULT 0,
    bills_introduced_current_term INTEGER DEFAULT 0,
    bills_passed_total INTEGER DEFAULT 0,
    bills_passed_current_term INTEGER DEFAULT 0,
    committee_attendance_rate DECIMAL(5,2), -- Percentage
    floor_speech_count INTEGER DEFAULT 0,
    amendments_proposed INTEGER DEFAULT 0,
    
    -- Voting Record
    total_votes_cast INTEGER DEFAULT 0,
    votes_yea INTEGER DEFAULT 0,
    votes_nay INTEGER DEFAULT 0,
    votes_present INTEGER DEFAULT 0,
    votes_absent INTEGER DEFAULT 0,
    voting_participation_rate DECIMAL(5,2), -- Percentage
    
    -- Transparency Score Components
    financial_disclosure_compliant BOOLEAN DEFAULT true,
    meeting_attendance_rate DECIMAL(5,2), -- Percentage
    public_communications_count INTEGER DEFAULT 0,
    public_communications_last_month INTEGER DEFAULT 0,
    transparency_score DECIMAL(5,2), -- 0-100 calculated score
    
    -- Comparative Metrics (vs city average)
    response_time_vs_avg DECIMAL(10,2), -- Percentage difference
    productivity_vs_avg DECIMAL(10,2), -- Percentage difference
    attendance_vs_avg DECIMAL(10,2), -- Percentage difference
    
    -- Metadata
    last_calculated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Ward comparison metrics (for ward-to-ward comparisons)
CREATE TABLE IF NOT EXISTS ward_metrics (
    id BIGSERIAL PRIMARY KEY,
    ward INTEGER UNIQUE NOT NULL,
    official_id INTEGER REFERENCES officials(id),
    
    -- Ward-specific metrics
    population INTEGER,
    median_income DECIMAL(10,2),
    unemployment_rate DECIMAL(5,2),
    crime_rate DECIMAL(10,2),
    
    -- Service delivery
    avg_service_request_resolution_days DECIMAL(10,2),
    total_service_requests INTEGER DEFAULT 0,
    service_requests_resolved INTEGER DEFAULT 0,
    
    -- Infrastructure
    pothole_repairs INTEGER DEFAULT 0,
    street_light_repairs INTEGER DEFAULT 0,
    tree_trimming_requests INTEGER DEFAULT 0,
    
    -- Community engagement
    community_meetings_held INTEGER DEFAULT 0,
    residents_served INTEGER DEFAULT 0,
    constituent_satisfaction_score DECIMAL(5,2),
    
    -- Metadata
    last_updated TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Legislative activity log (detailed tracking)
CREATE TABLE IF NOT EXISTS legislative_activity (
    id BIGSERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) NOT NULL,
    activity_type TEXT NOT NULL, -- 'bill_introduced', 'vote_cast', 'speech_given', 'committee_attendance'
    activity_date TIMESTAMP NOT NULL,
    matter_id TEXT, -- Reference to matters table
    description TEXT,
    metadata JSONB, -- Additional data specific to activity type
    created_at TIMESTAMP DEFAULT NOW()
);

-- Constituent service log
CREATE TABLE IF NOT EXISTS constituent_services (
    id BIGSERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) NOT NULL,
    service_type TEXT NOT NULL, -- 'case', 'office_hours', 'town_hall', 'communication'
    request_date TIMESTAMP NOT NULL,
    resolution_date TIMESTAMP,
    response_time_hours DECIMAL(10,2),
    status TEXT, -- 'open', 'resolved', 'closed'
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Transparency tracking
CREATE TABLE IF NOT EXISTS transparency_records (
    id BIGSERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) NOT NULL,
    record_type TEXT NOT NULL, -- 'financial_disclosure', 'meeting', 'public_statement'
    record_date TIMESTAMP NOT NULL,
    compliant BOOLEAN DEFAULT true,
    notes TEXT,
    source_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_official_metrics_official_id ON official_metrics(official_id);
CREATE INDEX IF NOT EXISTS idx_ward_metrics_ward ON ward_metrics(ward);
CREATE INDEX IF NOT EXISTS idx_ward_metrics_official_id ON ward_metrics(official_id);
CREATE INDEX IF NOT EXISTS idx_legislative_activity_official_id ON legislative_activity(official_id);
CREATE INDEX IF NOT EXISTS idx_legislative_activity_date ON legislative_activity(activity_date);
CREATE INDEX IF NOT EXISTS idx_legislative_activity_type ON legislative_activity(activity_type);
CREATE INDEX IF NOT EXISTS idx_constituent_services_official_id ON constituent_services(official_id);
CREATE INDEX IF NOT EXISTS idx_constituent_services_status ON constituent_services(status);
CREATE INDEX IF NOT EXISTS idx_transparency_records_official_id ON transparency_records(official_id);
CREATE INDEX IF NOT EXISTS idx_transparency_records_type ON transparency_records(record_type);

-- Update triggers
DROP TRIGGER IF EXISTS update_official_metrics_updated_at ON official_metrics;
CREATE TRIGGER update_official_metrics_updated_at 
    BEFORE UPDATE ON official_metrics 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_ward_metrics_updated_at ON ward_metrics;
CREATE TRIGGER update_ward_metrics_updated_at 
    BEFORE UPDATE ON ward_metrics 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_constituent_services_updated_at ON constituent_services;
CREATE TRIGGER update_constituent_services_updated_at 
    BEFORE UPDATE ON constituent_services 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- View for easy metrics access
CREATE OR REPLACE VIEW official_metrics_summary AS
SELECT 
    o.id,
    o.name,
    o.ward,
    o.party,
    o.role,
    om.avg_response_time_hours,
    om.cases_resolved_total,
    om.office_hours_held_total,
    om.town_halls_attended_total,
    om.bills_introduced_current_term,
    om.bills_passed_current_term,
    om.committee_attendance_rate,
    om.voting_participation_rate,
    om.transparency_score,
    om.response_time_vs_avg,
    om.productivity_vs_avg,
    om.attendance_vs_avg,
    om.last_calculated_at
FROM officials o
LEFT JOIN official_metrics om ON o.id = om.official_id;
