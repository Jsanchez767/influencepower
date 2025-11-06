-- =====================================================
-- MIGRATION: Drop Old Schema & Apply New Scalable Schema
-- =====================================================
-- This migration:
-- 1. Drops all existing tables
-- 2. Applies new position-based schema
-- 3. Adds multi-city/multi-government support
-- 4. Seeds initial data for Chicago
-- =====================================================

-- =====================================================
-- STEP 1: DROP OLD TABLES (Clean Slate)
-- =====================================================

-- Drop old views first
DROP VIEW IF EXISTS official_metrics_summary CASCADE;
DROP VIEW IF EXISTS current_officials CASCADE;
DROP VIEW IF EXISTS official_history CASCADE;

-- Drop old tables in reverse dependency order
DROP TABLE IF EXISTS transparency_records CASCADE;
DROP TABLE IF EXISTS constituent_services CASCADE;
DROP TABLE IF EXISTS legislative_activity CASCADE;
DROP TABLE IF EXISTS ward_metrics CASCADE;
DROP TABLE IF EXISTS official_metrics CASCADE;
DROP TABLE IF EXISTS person_metrics CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS vote_records CASCADE;
DROP TABLE IF EXISTS terms CASCADE;
DROP TABLE IF EXISTS people CASCADE;
DROP TABLE IF EXISTS positions CASCADE;
DROP TABLE IF EXISTS officials CASCADE;
DROP TABLE IF EXISTS matters CASCADE;

-- Drop old functions
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;

-- =====================================================
-- STEP 2: MULTI-JURISDICTION SUPPORT
-- =====================================================

-- Jurisdictions table (cities, states, federal)
CREATE TABLE IF NOT EXISTS jurisdictions (
  id SERIAL PRIMARY KEY,
  
  -- Identification
  name TEXT NOT NULL,                    -- 'Chicago', 'New York', 'Illinois', 'United States'
  jurisdiction_type TEXT NOT NULL,       -- 'city', 'county', 'state', 'federal'
  
  -- Location
  state_code TEXT,                       -- 'IL', 'NY', etc.
  country_code TEXT DEFAULT 'US',
  
  -- API Configuration
  api_type TEXT,                         -- 'legistar', 'open_states', 'congress_gov'
  api_base_url TEXT,                     -- 'https://webapi.legistar.com/v1/chicago'
  api_key TEXT,                          -- For APIs that require keys
  
  -- Settings
  timezone TEXT DEFAULT 'America/Chicago',
  population INTEGER,
  is_active BOOLEAN DEFAULT true,
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  UNIQUE(name, jurisdiction_type, state_code)
);

-- Index for lookups
CREATE INDEX IF NOT EXISTS idx_jurisdictions_type ON jurisdictions(jurisdiction_type);
CREATE INDEX IF NOT EXISTS idx_jurisdictions_active ON jurisdictions(is_active) WHERE is_active = true;

-- =====================================================
-- STEP 3: POSITIONS TABLE (Multi-Jurisdiction)
-- =====================================================

CREATE TABLE IF NOT EXISTS positions (
  id SERIAL PRIMARY KEY,
  
  -- Jurisdiction
  jurisdiction_id INTEGER NOT NULL REFERENCES jurisdictions(id) ON DELETE CASCADE,
  
  -- Position Details
  position_type TEXT NOT NULL,           -- 'alderman', 'mayor', 'senator', 'representative'
  district_number INTEGER,               -- Ward, district, etc. (NULL for at-large)
  district_name TEXT,                    -- 'Ward 1', '5th Congressional District'
  title TEXT NOT NULL,                   -- 'Alderman', 'Mayor', 'State Senator'
  
  -- Legislative Body
  body_name TEXT,                        -- 'City Council', 'State Senate', 'U.S. House'
  body_id INTEGER,                       -- From external API
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  UNIQUE(jurisdiction_id, position_type, district_number)
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_positions_jurisdiction ON positions(jurisdiction_id);
CREATE INDEX IF NOT EXISTS idx_positions_district ON positions(district_number);
CREATE INDEX IF NOT EXISTS idx_positions_type ON positions(position_type);
CREATE INDEX IF NOT EXISTS idx_positions_lookup ON positions(jurisdiction_id, position_type, district_number);

-- =====================================================
-- STEP 4: PEOPLE TABLE (Cross-Jurisdiction)
-- =====================================================

CREATE TABLE IF NOT EXISTS people (
  id SERIAL PRIMARY KEY,
  
  -- Identity
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  full_name TEXT NOT NULL,
  middle_name TEXT,
  suffix TEXT,                           -- 'Jr.', 'III', etc.
  
  -- External IDs (for syncing)
  external_ids JSONB DEFAULT '{}',       -- {"legistar_id": 162, "bioguide_id": "S000148"}
  
  -- Contact
  email TEXT,
  phone TEXT,
  website TEXT,
  
  -- Social Media
  twitter_handle TEXT,
  facebook_url TEXT,
  instagram_handle TEXT,
  
  -- Media
  image_url TEXT,
  headshot_last_updated TIMESTAMPTZ,
  
  -- Biographical
  date_of_birth DATE,
  party_affiliation TEXT,                -- 'Democrat', 'Republican', 'Independent'
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_people_full_name ON people(full_name);
CREATE INDEX IF NOT EXISTS idx_people_external_ids ON people USING GIN(external_ids);
CREATE INDEX IF NOT EXISTS idx_people_party ON people(party_affiliation);

-- =====================================================
-- STEP 5: TERMS TABLE (Who Holds What When)
-- =====================================================

CREATE TABLE IF NOT EXISTS terms (
  id SERIAL PRIMARY KEY,
  
  -- Relationships
  position_id INTEGER NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
  person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  
  -- Term Dates
  start_date DATE NOT NULL,
  end_date DATE,                         -- NULL = current term
  
  -- External Reference
  external_id INTEGER,                   -- OfficeRecordId, etc.
  external_guid TEXT,
  
  -- Term Details
  term_number INTEGER,                   -- 1st term, 2nd term, etc.
  election_type TEXT,                    -- 'general', 'special', 'appointment'
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  CHECK (end_date IS NULL OR end_date >= start_date),
  UNIQUE(position_id, person_id, start_date)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_terms_position ON terms(position_id);
CREATE INDEX IF NOT EXISTS idx_terms_person ON terms(person_id);
CREATE INDEX IF NOT EXISTS idx_terms_current ON terms(end_date) WHERE end_date IS NULL;
CREATE INDEX IF NOT EXISTS idx_terms_dates ON terms(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_terms_external_id ON terms(external_id);

-- Create a function to check if a term is current
CREATE OR REPLACE FUNCTION is_term_current(term_end_date DATE) 
RETURNS BOOLEAN AS $$
BEGIN
  RETURN term_end_date IS NULL OR term_end_date > CURRENT_DATE;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- =====================================================
-- STEP 6: LEGISLATION (Bills, Ordinances, Resolutions)
-- =====================================================

CREATE TABLE IF NOT EXISTS legislation (
  id SERIAL PRIMARY KEY,
  
  -- Jurisdiction
  jurisdiction_id INTEGER NOT NULL REFERENCES jurisdictions(id) ON DELETE CASCADE,
  
  -- Identification
  external_id INTEGER,                   -- MatterId from API
  external_guid TEXT,
  file_number TEXT,                      -- 'O2023-1234', 'HR 1234'
  
  -- Content
  title TEXT NOT NULL,
  summary TEXT,
  full_text TEXT,
  
  -- Type & Status
  legislation_type TEXT,                 -- 'ordinance', 'resolution', 'bill'
  status TEXT,                           -- 'introduced', 'passed', 'failed', 'vetoed'
  
  -- Dates
  introduced_date DATE,
  passed_date DATE,
  enacted_date DATE,
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_legislation_jurisdiction ON legislation(jurisdiction_id);
CREATE INDEX IF NOT EXISTS idx_legislation_external_id ON legislation(external_id);
CREATE INDEX IF NOT EXISTS idx_legislation_status ON legislation(status);
CREATE INDEX IF NOT EXISTS idx_legislation_dates ON legislation(introduced_date, passed_date);

-- =====================================================
-- STEP 7: SPONSORSHIPS (Who Introduced What)
-- =====================================================

CREATE TABLE IF NOT EXISTS sponsorships (
  id SERIAL PRIMARY KEY,
  
  legislation_id INTEGER NOT NULL REFERENCES legislation(id) ON DELETE CASCADE,
  person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  
  sponsorship_type TEXT NOT NULL,        -- 'primary', 'co-sponsor'
  position INTEGER,                      -- Order of sponsorship
  
  created_at TIMESTAMPTZ DEFAULT NOW(),
  
  UNIQUE(legislation_id, person_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_sponsorships_legislation ON sponsorships(legislation_id);
CREATE INDEX IF NOT EXISTS idx_sponsorships_person ON sponsorships(person_id);
CREATE INDEX IF NOT EXISTS idx_sponsorships_type ON sponsorships(sponsorship_type);

-- =====================================================
-- STEP 8: VOTES (Voting Records)
-- =====================================================

CREATE TABLE IF NOT EXISTS votes (
  id SERIAL PRIMARY KEY,
  
  -- Who voted
  person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  
  -- What they voted on
  legislation_id INTEGER REFERENCES legislation(id) ON DELETE CASCADE,
  
  -- External IDs (for syncing)
  external_vote_id INTEGER,
  external_event_id INTEGER,             -- Meeting/session ID
  
  -- The vote
  vote_value TEXT NOT NULL,              -- 'yea', 'nay', 'abstain', 'absent', 'present'
  vote_date DATE NOT NULL,
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  UNIQUE(person_id, legislation_id, vote_date)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_votes_person ON votes(person_id);
CREATE INDEX IF NOT EXISTS idx_votes_legislation ON votes(legislation_id);
CREATE INDEX IF NOT EXISTS idx_votes_date ON votes(vote_date DESC);
CREATE INDEX IF NOT EXISTS idx_votes_value ON votes(vote_value);
CREATE INDEX IF NOT EXISTS idx_votes_external ON votes(external_vote_id);

-- =====================================================
-- STEP 9: METRICS (Denormalized for Performance)
-- =====================================================

CREATE TABLE IF NOT EXISTS person_metrics (
  person_id INTEGER PRIMARY KEY REFERENCES people(id) ON DELETE CASCADE,
  
  -- Legislative Activity
  bills_introduced INTEGER DEFAULT 0,
  bills_passed INTEGER DEFAULT 0,
  bills_failed INTEGER DEFAULT 0,
  amendments_proposed INTEGER DEFAULT 0,
  
  -- Voting Record
  total_votes INTEGER DEFAULT 0,
  votes_yea INTEGER DEFAULT 0,
  votes_nay INTEGER DEFAULT 0,
  votes_abstain INTEGER DEFAULT 0,
  votes_absent INTEGER DEFAULT 0,
  attendance_rate DECIMAL(5,2) DEFAULT 0,
  
  -- Scores (0-100)
  legislative_impact_score DECIMAL(5,2) DEFAULT 0,
  constituent_engagement_score DECIMAL(5,2) DEFAULT 0,
  transparency_score DECIMAL(5,2) DEFAULT 0,
  overall_score DECIMAL(5,2) DEFAULT 0,
  
  -- Timestamps
  last_calculated_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- =====================================================
-- STEP 10: VIEWS FOR EASY QUERIES
-- =====================================================

-- Current officials across all jurisdictions
CREATE OR REPLACE VIEW current_officials AS
SELECT 
  j.id as jurisdiction_id,
  j.name as jurisdiction_name,
  j.jurisdiction_type,
  pos.id as position_id,
  pos.position_type,
  pos.district_number,
  pos.district_name,
  pos.title,
  p.id as person_id,
  p.full_name,
  p.first_name,
  p.last_name,
  p.party_affiliation,
  p.email,
  p.phone,
  p.website,
  p.image_url,
  t.start_date as term_start,
  t.end_date as term_end,
  t.term_number,
  m.overall_score,
  m.legislative_impact_score,
  m.constituent_engagement_score,
  m.transparency_score,
  m.attendance_rate
FROM jurisdictions j
JOIN positions pos ON pos.jurisdiction_id = j.id
JOIN terms t ON t.position_id = pos.id
JOIN people p ON p.id = t.person_id
LEFT JOIN person_metrics m ON m.person_id = p.id
WHERE (t.end_date IS NULL OR t.end_date > CURRENT_DATE)
  AND j.is_active = true
ORDER BY j.name, pos.district_number NULLS FIRST, pos.position_type;

-- Historical terms
CREATE OR REPLACE VIEW term_history AS
SELECT 
  j.name as jurisdiction_name,
  pos.district_number,
  pos.title,
  p.full_name,
  p.party_affiliation,
  t.start_date,
  t.end_date,
  EXTRACT(YEAR FROM AGE(COALESCE(t.end_date, CURRENT_DATE), t.start_date)) as years_served,
  t.term_number
FROM jurisdictions j
JOIN positions pos ON pos.jurisdiction_id = j.id
JOIN terms t ON t.position_id = pos.id
JOIN people p ON p.id = t.person_id
ORDER BY j.name, pos.district_number NULLS FIRST, t.start_date DESC;

-- =====================================================
-- STEP 11: TRIGGERS
-- =====================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at
DROP TRIGGER IF EXISTS update_jurisdictions_updated_at ON jurisdictions;
CREATE TRIGGER update_jurisdictions_updated_at BEFORE UPDATE ON jurisdictions
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_positions_updated_at ON positions;
CREATE TRIGGER update_positions_updated_at BEFORE UPDATE ON positions
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_people_updated_at ON people;
CREATE TRIGGER update_people_updated_at BEFORE UPDATE ON people
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_terms_updated_at ON terms;
CREATE TRIGGER update_terms_updated_at BEFORE UPDATE ON terms
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_legislation_updated_at ON legislation;
CREATE TRIGGER update_legislation_updated_at BEFORE UPDATE ON legislation
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_person_metrics_updated_at ON person_metrics;
CREATE TRIGGER update_person_metrics_updated_at BEFORE UPDATE ON person_metrics
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- STEP 12: SEED DATA - Chicago
-- =====================================================

-- Add Chicago jurisdiction
INSERT INTO jurisdictions (name, jurisdiction_type, state_code, api_type, api_base_url, population)
VALUES (
  'Chicago',
  'city',
  'IL',
  'legistar',
  'https://webapi.legistar.com/v1/chicago',
  2700000
) ON CONFLICT (name, jurisdiction_type, state_code) DO NOTHING;

-- Get Chicago's ID and create positions
DO $$
DECLARE
  chicago_id INTEGER;
BEGIN
  SELECT id INTO chicago_id FROM jurisdictions WHERE name = 'Chicago' AND jurisdiction_type = 'city';
  
  -- Create all 50 ward alderman positions
  INSERT INTO positions (jurisdiction_id, position_type, district_number, district_name, title, body_name, body_id)
  SELECT 
    chicago_id,
    'alderman',
    generate_series(1, 50),
    'Ward ' || generate_series(1, 50),
    'Alderman',
    'City Council',
    138
  ON CONFLICT (jurisdiction_id, position_type, district_number) DO NOTHING;
  
  -- Create citywide positions
  INSERT INTO positions (jurisdiction_id, position_type, district_number, district_name, title, body_name)
  VALUES
    (chicago_id, 'mayor', NULL, 'Citywide', 'Mayor', 'City Council'),
    (chicago_id, 'clerk', NULL, 'Citywide', 'City Clerk', NULL),
    (chicago_id, 'treasurer', NULL, 'Citywide', 'City Treasurer', NULL)
  ON CONFLICT (jurisdiction_id, position_type, district_number) DO NOTHING;
END $$;

-- =====================================================
-- MIGRATION COMPLETE! ✅
-- =====================================================
-- Next steps:
-- 1. Run: go run scripts/sync_positions_v2.go
-- 2. This will populate people and terms from City API
-- 3. Then run: go run scripts/calculate_metrics.go
-- =====================================================
