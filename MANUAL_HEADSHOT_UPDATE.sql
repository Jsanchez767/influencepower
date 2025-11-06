-- First, add the column if needed
ALTER TABLE officials ADD COLUMN IF NOT EXISTS city_api_person_id TEXT;

-- Update headshots and city_api_person_id for known officials
-- You'll need to get these from the City API or manually map them

-- Example updates (you'll need to fill in the actual person_ids and image URLs from City API):
-- UPDATE officials SET 
--   city_api_person_id = 'person-guid-here',
--   image_url = 'https://chicago.legistar.com/people/person-image.jpg'
-- WHERE name = 'Daniel La Spata';

-- For now, let's at least ensure the column exists
CREATE INDEX IF NOT EXISTS idx_officials_city_api_person_id ON officials(city_api_person_id);
