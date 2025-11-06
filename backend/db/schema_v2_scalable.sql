-- =====================================================
-- IMPROVED HORIZONTALLY SCALABLE SCHEMA
-- =====================================================
-- This schema separates concerns and uses position-based
-- lookups instead of name matching, making it future-proof
-- when officials change.

-- =====================================================
-- 1. POSITIONS TABLE (The Source of Truth)
-- =====================================================
-- Positions are stable (Ward 1 Alderman will always exist)
-- People come and go, but positions remain
CREATE TABLE IF NOT EXISTS positions (
  id SERIAL PRIMARY KEY,
  position_type TEXT NOT NULL, -- 'alderman', 'mayor', 'clerk', 'treasurer'
  ward INTEGER, -- NULL for citywide positions
  title TEXT NOT NULL, -- 'Alderman', 'Mayor', 'City Clerk', etc.
  body_name TEXT, -- 'City Council', etc.
  body_id INTEGER, -- From City API
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  UNIQUE(position_type, ward), -- Only one alderman per ward
  CHECK (
    (position_type = 'alderman' AND ward IS NOT NULL) OR
    (position_type IN ('mayor', 'clerk', 'treasurer') AND ward IS NULL)
  )
);

-- Index for fast ward lookups
CREATE INDEX IF NOT EXISTS idx_positions_ward ON positions(ward);
CREATE INDEX IF NOT EXISTS idx_positions_type ON positions(position_type);

-- =====================================================
-- 2. PEOPLE TABLE (Officials as Individuals)
-- =====================================================
-- Stores information about people, separate from their positions
CREATE TABLE IF NOT EXISTS people (
  id SERIAL PRIMARY KEY,
  
  -- Identity (from City API)
  city_api_person_id INTEGER UNIQUE, -- PersonId from Legistar
  city_api_person_guid TEXT UNIQUE, -- PersonGuid from Legistar
  
  -- Name
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  full_name TEXT NOT NULL,
  
  -- Contact
  email TEXT,
  phone TEXT,
  website TEXT,
  
  -- Media
  image_url TEXT,
  headshot_last_updated TIMESTAMPTZ,
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for lookups
CREATE INDEX IF NOT EXISTS idx_people_city_api_person_id ON people(city_api_person_id);
CREATE INDEX IF NOT EXISTS idx_people_full_name ON people(full_name);

-- =====================================================
-- 3. TERMS TABLE (Who Holds What Position When)
-- =====================================================
-- Links people to positions over time
-- This is the key to handling position changes!
CREATE TABLE IF NOT EXISTS terms (
  id SERIAL PRIMARY KEY,
  
  -- Relationships
  position_id INTEGER NOT NULL REFERENCES positions(id) ON DELETE CASCADE,
  person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  
  -- Term dates (from City API OfficeRecords)
  start_date DATE NOT NULL,
  end_date DATE, -- NULL means current term
  
  -- City API reference
  city_api_office_record_id INTEGER UNIQUE,
  city_api_office_record_guid TEXT UNIQUE,
  
  -- Metadata
  is_current BOOLEAN GENERATED ALWAYS AS (end_date IS NULL OR end_date > CURRENT_DATE) STORED,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  CHECK (end_date IS NULL OR end_date >= start_date)
);

-- Indexes for fast queries
CREATE INDEX IF NOT EXISTS idx_terms_position_id ON terms(position_id);
CREATE INDEX IF NOT EXISTS idx_terms_person_id ON terms(person_id);
CREATE INDEX IF NOT EXISTS idx_terms_current ON terms(is_current) WHERE is_current = true;
CREATE INDEX IF NOT EXISTS idx_terms_dates ON terms(start_date, end_date);

-- =====================================================
-- 4. VOTING RECORDS (Tied to People, Not Positions)
-- =====================================================
-- Votes are cast by people, not positions
CREATE TABLE IF NOT EXISTS votes (
  id SERIAL PRIMARY KEY,
  
  -- Who voted
  person_id INTEGER NOT NULL REFERENCES people(id) ON DELETE CASCADE,
  
  -- What they voted on
  matter_id INTEGER NOT NULL, -- From City API
  matter_file TEXT,
  matter_title TEXT,
  
  -- The vote
  vote_value TEXT NOT NULL, -- 'Yea', 'Nay', 'Abstain', etc.
  vote_date DATE NOT NULL,
  
  -- City API reference
  city_api_vote_id INTEGER UNIQUE,
  city_api_vote_guid TEXT,
  
  -- Event context
  event_id INTEGER, -- From City API
  
  -- Metadata
  created_at TIMESTAMPTZ DEFAULT NOW(),
  
  -- Constraints
  UNIQUE(person_id, matter_id, vote_date)
);

-- Indexes for analytics
CREATE INDEX IF NOT EXISTS idx_votes_person_id ON votes(person_id);
CREATE INDEX IF NOT EXISTS idx_votes_matter_id ON votes(matter_id);
CREATE INDEX IF NOT EXISTS idx_votes_date ON votes(vote_date DESC);
CREATE INDEX IF NOT EXISTS idx_votes_value ON votes(vote_value);

-- =====================================================
-- 5. METRICS (Denormalized for Performance)
-- =====================================================
-- Keep existing metrics structure but reference people
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
-- 6. VIEWS FOR EASY QUERIES
-- =====================================================

-- Current officials with all their info
CREATE OR REPLACE VIEW current_officials AS
SELECT 
  pos.id as position_id,
  pos.position_type,
  pos.ward,
  pos.title,
  p.id as person_id,
  p.city_api_person_id,
  p.full_name,
  p.first_name,
  p.last_name,
  p.email,
  p.phone,
  p.website,
  p.image_url,
  t.start_date as term_start,
  t.end_date as term_end,
  m.overall_score,
  m.legislative_impact_score,
  m.constituent_engagement_score,
  m.transparency_score
FROM positions pos
JOIN terms t ON t.position_id = pos.id
JOIN people p ON p.id = t.person_id
LEFT JOIN person_metrics m ON m.person_id = p.id
WHERE t.is_current = true
ORDER BY pos.ward NULLS FIRST, pos.position_type;

-- Historical officials (for tracking changes)
CREATE OR REPLACE VIEW official_history AS
SELECT 
  pos.ward,
  pos.title,
  p.full_name,
  p.email,
  t.start_date,
  t.end_date,
  EXTRACT(YEAR FROM AGE(COALESCE(t.end_date, CURRENT_DATE), t.start_date)) as years_served
FROM positions pos
JOIN terms t ON t.position_id = pos.id
JOIN people p ON p.id = t.person_id
ORDER BY pos.ward NULLS FIRST, t.start_date DESC;

-- =====================================================
-- 7. TRIGGERS FOR AUTO-UPDATES
-- =====================================================

-- Update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_positions_updated_at ON positions;
CREATE TRIGGER update_positions_updated_at BEFORE UPDATE ON positions
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_people_updated_at ON people;
CREATE TRIGGER update_people_updated_at BEFORE UPDATE ON people
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_terms_updated_at ON terms;
CREATE TRIGGER update_terms_updated_at BEFORE UPDATE ON terms
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- 8. SAMPLE DATA - SEED POSITIONS
-- =====================================================

-- Insert all 50 ward alderman positions
INSERT INTO positions (position_type, ward, title, body_name, body_id)
SELECT 
  'alderman',
  generate_series(1, 50),
  'Alderman',
  'City Council',
  138
ON CONFLICT (position_type, ward) DO NOTHING;

-- Insert citywide positions
INSERT INTO positions (position_type, ward, title, body_name) VALUES
  ('mayor', NULL, 'Mayor', 'City Council'),
  ('clerk', NULL, 'City Clerk', NULL),
  ('treasurer', NULL, 'City Treasurer', NULL)
ON CONFLICT (position_type, ward) DO NOTHING;

-- =====================================================
-- BENEFITS OF THIS SCHEMA:
-- =====================================================
-- ✅ Position-based lookups (not name-based)
-- ✅ Historical tracking of who held positions
-- ✅ Handles official changes automatically
-- ✅ Normalized data (people separate from positions)
-- ✅ Horizontally scalable (can shard by ward)
-- ✅ Fast queries with proper indexes
-- ✅ Current vs historical data separated
-- ✅ Votes tied to people, not positions
-- ✅ Easy to add new position types
-- ✅ City API IDs for reliable syncing
