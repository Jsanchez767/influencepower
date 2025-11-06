# Chicago City Clerk API Integration

## Overview

This document outlines our **Hybrid Approach** for integrating the Chicago City Clerk ELMS API with our existing Supabase database.

## Architecture: Hybrid Sync + On-Demand

We use a **scheduled sync** for most data (legislation, votes, events) while keeping official profiles in Supabase as the source of truth.

### Why Hybrid?

1. **Data Ownership**: We control official profiles (photos, bios, social links)
2. **Performance**: Cached legislative data for fast queries
3. **Freshness**: Recent votes and events updated regularly
4. **Reliability**: App works even if City API is down

## Data Flow

```
┌─────────────────────┐
│  Chicago City Clerk │
│   ELMS API          │
│  (Legistar v1)      │
└──────────┬──────────┘
           │
           │ Scheduled Sync (Daily)
           ▼
┌─────────────────────┐
│   Supabase          │
│   PostgreSQL        │
│                     │
│  ├─ officials       │ ← Our data (source of truth)
│  ├─ matters         │ ← Synced from City API
│  ├─ votes           │ ← Synced from City API
│  ├─ events          │ ← Synced from City API
│  ├─ bodies          │ ← Synced from City API
│  └─ event_items     │ ← Synced from City API
└──────────┬──────────┘
           │
           │ REST API
           ▼
┌─────────────────────┐
│   Go Backend        │
│   (Render)          │
└──────────┬──────────┘
           │
           │ HTTPS
           ▼
┌─────────────────────┐
│  React Frontend     │
│   (Vercel)          │
└─────────────────────┘
```

## Implementation

### 1. Database Schema (`backend/db/schema_city_api.sql`)

Extended tables to store City API data:

- **matters**: Legislation (ordinances, resolutions, appointments)
- **votes**: Individual votes on matters
- **events**: Council meetings and hearings
- **event_items**: Agenda items linked to matters
- **bodies**: Committees and council bodies

All tables have:
- Unique constraints on API IDs
- Foreign keys to link data
- Indexes for performance
- Timestamps for tracking

### 2. City API Client (`backend/cityapi/client.go`)

Go client for the Chicago City Clerk ELMS API:

```go
client := cityapi.NewClient()

// Fetch recent legislation
matters, err := client.GetMatters(map[string]string{
    "$filter": "MatterIntroDate ge datetime'2024-01-01'",
    "$top": "100",
})

// Fetch votes for a specific matter
votes, err := client.GetVotes(matterID)

// Fetch upcoming events
events, err := client.GetEvents(params)
```

**Supported Endpoints**:
- `GET /matters` - Legislation
- `GET /matters/{id}` - Specific matter
- `GET /matters/{id}/votes` - Votes on matter
- `GET /events` - Meetings/hearings
- `GET /events/{id}` - Specific event
- `GET /bodies` - Committees
- `GET /persons` - Officials (for reference)

### 3. Sync Script (`backend/scripts/sync_city_api.go`)

Periodic sync job that:

1. **Syncs Bodies/Committees**
   - Fetches all committees from City API
   - Upserts into Supabase `bodies` table

2. **Syncs Recent Matters** (last 6 months)
   - Fetches legislation with filters
   - Stores sponsors, attachments as JSONB
   - For each matter, fetches and stores votes

3. **Syncs Recent Events** (last 3 months)
   - Fetches council meetings
   - Stores agenda items and links to matters

**Rate Limiting**: 100ms delay between requests

**Limits**:
- Initial sync: 100 matters, 50 events
- Can be increased for production

## Running the Sync

### Setup

1. **Apply Database Schema**:
   ```bash
   # In Supabase SQL Editor
   cat backend/db/schema_city_api.sql
   # Copy and run
   ```

2. **Set Environment Variables**:
   ```bash
   export SUPABASE_URL="https://your-project.supabase.co/rest/v1"
   export SUPABASE_SERVICE_ROLE_KEY="your-service-role-key"
   ```

3. **Run Sync Script**:
   ```bash
   cd backend
   go run scripts/sync_city_api.go
   ```

### Expected Output

```
🔄 Starting sync with Chicago City Clerk API...
📋 Syncing bodies/committees...
Found 25 bodies
  ✓ Synced: City Council
  ✓ Synced: Committee on Finance
  ...
👥 Syncing persons...
Found 53 persons
📄 Syncing recent matters...
Found 450 recent matters
  ✓ Synced: O2024-1234 - Zoning Amendment for Ward 1
    📊 Found 51 votes for matter 12345
  ...
📅 Syncing recent events...
Found 18 recent events
  ✓ Synced: City Council on 2024-01-17
  ...
✅ Sync complete!
```

## Scheduled Sync (Production)

### Option 1: Cron Job on Render

Add to `render.yaml`:
```yaml
- type: cron
  name: city-api-sync
  schedule: "0 2 * * *" # Daily at 2 AM
  buildCommand: cd backend && go build -o bin/sync scripts/sync_city_api.go
  startCommand: cd backend && ./bin/sync
```

### Option 2: GitHub Actions

```yaml
name: Sync City API
on:
  schedule:
    - cron: '0 2 * * *' # Daily at 2 AM
  workflow_dispatch: # Manual trigger

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run sync
        env:
          SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
          SUPABASE_SERVICE_ROLE_KEY: ${{ secrets.SUPABASE_SERVICE_ROLE_KEY }}
        run: |
          cd backend
          go run scripts/sync_city_api.go
```

### Option 3: Supabase Edge Function (Scheduled)

Create a scheduled function that triggers the sync endpoint.

## API Endpoints (Backend)

Add new handlers to expose synced data:

```go
// Get recent legislation
GET /api/v1/legislation
GET /api/v1/legislation/{id}

// Get votes for an official
GET /api/v1/officials/{id}/votes

// Get upcoming events
GET /api/v1/events

// Get legislation by committee
GET /api/v1/committees/{id}/legislation
```

## Frontend Integration

### Display Voting History

```typescript
// In official profile modal
const [votingHistory, setVotingHistory] = useState([]);

useEffect(() => {
  fetch(`${API_URL}/officials/${official.id}/votes`)
    .then(res => res.json())
    .then(data => setVotingHistory(data));
}, [official.id]);

// Render
{votingHistory.map(vote => (
  <div key={vote.id}>
    <h4>{vote.matter_title}</h4>
    <span className={vote.vote_value === 'Yea' ? 'text-green' : 'text-red'}>
      {vote.vote_value}
    </span>
    <time>{new Date(vote.vote_date).toLocaleDateString()}</time>
  </div>
))}
```

### Display Recent Legislation

```typescript
const [legislation, setLegislation] = useState([]);

useEffect(() => {
  fetch(`${API_URL}/legislation?limit=10`)
    .then(res => res.json())
    .then(data => setLegislation(data));
}, []);
```

## Data Matching

The City API has `PersonID` that we need to match to our `officials.id`.

### Manual Mapping Table

Create a mapping table:

```sql
CREATE TABLE official_person_mapping (
    official_id INTEGER REFERENCES officials(id),
    person_id INTEGER,
    person_name TEXT,
    verified BOOLEAN DEFAULT false,
    PRIMARY KEY (official_id, person_id)
);
```

Then populate by name matching:

```sql
-- Match by name similarity
INSERT INTO official_person_mapping (official_id, person_id, person_name)
SELECT 
    o.id,
    v.person_id,
    v.person_name
FROM officials o
JOIN votes v ON LOWER(o.name) = LOWER(v.person_name)
GROUP BY o.id, v.person_id, v.person_name;
```

## Benefits of This Approach

✅ **Fast Queries**: Data cached in Supabase  
✅ **Enriched Data**: Combine City API + our custom data  
✅ **Reliability**: App works offline from City API  
✅ **Historical Data**: Keep full voting history  
✅ **Custom Features**: Add fields City API doesn't have  
✅ **Search**: Full-text search on legislation  
✅ **Analytics**: Query patterns, trends, statistics  

## Next Steps

1. ✅ Create database schema
2. ✅ Create City API client
3. ✅ Create sync script
4. ⏳ Apply schema to Supabase
5. ⏳ Run initial sync
6. ⏳ Add backend endpoints
7. ⏳ Update frontend UI
8. ⏳ Set up scheduled sync
9. ⏳ Create person-to-official mapping

## Resources

- [Chicago City Clerk API Documentation](https://api.chicityclerkelms.chicago.gov)
- [Legistar API Reference](https://webapi.legistar.com)
- Supabase Project: https://datckstcdyldpbvfxykw.supabase.co
- Backend API: https://backend.maticsapp.com/api
