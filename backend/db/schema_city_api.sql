-- Extended schema for Chicago City Clerk ELMS API integration

-- Legislation/Matters table
CREATE TABLE IF NOT EXISTS matters (
    id BIGSERIAL PRIMARY KEY,
    matter_id TEXT UNIQUE NOT NULL,
    matter_file TEXT,
    matter_name TEXT,
    matter_title TEXT,
    matter_type_id INTEGER,
    matter_type_name TEXT,
    matter_status_id INTEGER,
    matter_status_name TEXT,
    matter_intro_date TIMESTAMP,
    matter_agenda_date TIMESTAMP,
    matter_passed_date TIMESTAMP,
    matter_enactment_date TIMESTAMP,
    matter_enactment_number TEXT,
    matter_requester TEXT,
    matter_sponsors JSONB,
    matter_attachments JSONB,
    matter_text TEXT,
    matter_version TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Votes table (links to matters and officials)
CREATE TABLE IF NOT EXISTS votes (
    id BIGSERIAL PRIMARY KEY,
    vote_id TEXT UNIQUE,
    matter_id TEXT REFERENCES matters(matter_id),
    official_id INTEGER REFERENCES officials(id),
    person_id INTEGER, -- from city API
    person_name TEXT,
    vote_value TEXT, -- 'Yea', 'Nay', 'Abstain', 'Present', etc.
    vote_date TIMESTAMP,
    vote_event_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Events/Meetings table
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    event_id TEXT UNIQUE NOT NULL,
    event_body_id INTEGER,
    event_body_name TEXT,
    event_date TIMESTAMP,
    event_time TEXT,
    event_location TEXT,
    event_agenda_file TEXT,
    event_minutes_file TEXT,
    event_video_url TEXT,
    event_items JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Event Items (agenda items linked to matters and events)
CREATE TABLE IF NOT EXISTS event_items (
    id BIGSERIAL PRIMARY KEY,
    event_item_id TEXT UNIQUE NOT NULL,
    event_id TEXT REFERENCES events(event_id),
    matter_id TEXT REFERENCES matters(matter_id),
    item_agenda_sequence INTEGER,
    item_agenda_number TEXT,
    item_action TEXT,
    item_action_text TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Body/Committee types from API
CREATE TABLE IF NOT EXISTS bodies (
    id BIGSERIAL PRIMARY KEY,
    body_id INTEGER UNIQUE NOT NULL,
    body_name TEXT,
    body_type_id INTEGER,
    body_type_name TEXT,
    body_meet_flag INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_matters_matter_id ON matters(matter_id);
CREATE INDEX IF NOT EXISTS idx_matters_status ON matters(matter_status_name);
CREATE INDEX IF NOT EXISTS idx_matters_intro_date ON matters(matter_intro_date);
CREATE INDEX IF NOT EXISTS idx_votes_matter_id ON votes(matter_id);
CREATE INDEX IF NOT EXISTS idx_votes_official_id ON votes(official_id);
CREATE INDEX IF NOT EXISTS idx_votes_person_id ON votes(person_id);
CREATE INDEX IF NOT EXISTS idx_events_event_id ON events(event_id);
CREATE INDEX IF NOT EXISTS idx_events_date ON events(event_date);
CREATE INDEX IF NOT EXISTS idx_event_items_event_id ON event_items(event_id);
CREATE INDEX IF NOT EXISTS idx_event_items_matter_id ON event_items(matter_id);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_matters_updated_at BEFORE UPDATE ON matters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_votes_updated_at BEFORE UPDATE ON votes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_event_items_updated_at BEFORE UPDATE ON event_items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bodies_updated_at BEFORE UPDATE ON bodies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
