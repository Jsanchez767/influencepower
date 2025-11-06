-- Add city_api_person_id column to officials table if it doesn't exist
ALTER TABLE officials ADD COLUMN IF NOT EXISTS city_api_person_id TEXT;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_officials_city_api_person_id ON officials(city_api_person_id);
