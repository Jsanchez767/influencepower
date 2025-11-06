-- Create officials table
CREATE TABLE IF NOT EXISTS officials (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ward INTEGER,
    party VARCHAR(50) NOT NULL,
    role VARCHAR(100) NOT NULL,
    contact VARCHAR(255),
    email VARCHAR(255),
    image_url TEXT,
    city_api_person_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create voting_records table
CREATE TABLE IF NOT EXISTS voting_records (
    id SERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) ON DELETE CASCADE,
    bill_title VARCHAR(500) NOT NULL,
    vote VARCHAR(20) NOT NULL CHECK (vote IN ('yes', 'no', 'abstain')),
    vote_date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create committees table
CREATE TABLE IF NOT EXISTS committees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create official_committees junction table
CREATE TABLE IF NOT EXISTS official_committees (
    id SERIAL PRIMARY KEY,
    official_id INTEGER REFERENCES officials(id) ON DELETE CASCADE,
    committee_id INTEGER REFERENCES committees(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member' CHECK (role IN ('member', 'chair', 'vice-chair')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(official_id, committee_id)
);

-- Create ward_statistics table
CREATE TABLE IF NOT EXISTS ward_statistics (
    id SERIAL PRIMARY KEY,
    ward INTEGER NOT NULL UNIQUE,
    infrastructure_spending DECIMAL(12, 2) DEFAULT 0,
    response_311_time DECIMAL(5, 2) DEFAULT 0,
    affordable_housing INTEGER DEFAULT 0,
    permit_processing DECIMAL(5, 2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_officials_ward ON officials(ward);
CREATE INDEX IF NOT EXISTS idx_officials_party ON officials(party);
CREATE INDEX IF NOT EXISTS idx_voting_records_official ON voting_records(official_id);
CREATE INDEX IF NOT EXISTS idx_voting_records_date ON voting_records(vote_date);
CREATE INDEX IF NOT EXISTS idx_official_committees_official ON official_committees(official_id);
CREATE INDEX IF NOT EXISTS idx_official_committees_committee ON official_committees(committee_id);
CREATE INDEX IF NOT EXISTS idx_ward_statistics_ward ON ward_statistics(ward);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_officials_updated_at BEFORE UPDATE ON officials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ward_statistics_updated_at BEFORE UPDATE ON ward_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
