# InfluencePower Backend

Go backend API for the InfluencePower Chicago City Council visualization platform.

## Features

- RESTful API built with Go
- Supabase (PostgreSQL) database integration
- CRUD operations for officials, voting records, committees, and ward statistics
- CORS enabled for frontend integration

## Prerequisites

- Go 1.21 or higher
- Supabase account with a project set up

## Setup

1. **Clone the repository**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Copy `.env.example` to `.env` and fill in your Supabase credentials:
   ```bash
   cp .env.example .env
   ```

4. **Set up the database**
   
   Run the SQL schema in your Supabase SQL editor:
   ```bash
   # The schema is located in db/schema.sql
   ```
   
   Go to your Supabase project → SQL Editor → New Query → Paste the contents of `db/schema.sql` → Run

5. **Run the server**
   ```bash
   go run main.go
   ```
   
   The server will start on `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check API health status

### Officials
- `GET /api/v1/officials` - Get all officials
- `GET /api/v1/officials/{id}` - Get official by ID
- `POST /api/v1/officials` - Create new official
- `PUT /api/v1/officials/{id}` - Update official
- `DELETE /api/v1/officials/{id}` - Delete official
- `GET /api/v1/officials/party/{party}` - Get officials by party (democrat/republican)
- `GET /api/v1/officials/ward/{ward}` - Get officials by ward number

### Voting Records
- `GET /api/v1/officials/{id}/voting-records` - Get voting records for an official
- `POST /api/v1/voting-records` - Create new voting record

### Ward Statistics
- `GET /api/v1/wards/{ward}/statistics` - Get statistics for a specific ward

### Committees
- `GET /api/v1/committees` - Get all committees
- `GET /api/v1/officials/{id}/committees` - Get committees for an official

## Project Structure

```
backend/
├── main.go                 # Entry point
├── go.mod                  # Go dependencies
├── .env                    # Environment variables (not in git)
├── .env.example            # Environment variables template
├── api/
│   └── routes.go          # API route definitions
├── db/
│   ├── supabase.go        # Supabase client initialization
│   └── schema.sql         # Database schema
├── handlers/
│   └── handlers.go        # HTTP request handlers
└── models/
    └── models.go          # Data models
```

## Database Schema

### Tables
- **officials** - City officials (Mayor, Aldermen, City Clerk)
- **voting_records** - Voting history for each official
- **committees** - City committees
- **official_committees** - Junction table for official-committee relationships
- **ward_statistics** - Statistical data for each ward

## Development

To add new endpoints:
1. Define the model in `models/models.go`
2. Create handler functions in `handlers/handlers.go`
3. Register routes in `api/routes.go`

## Deployment

This backend can be deployed to:
- Heroku
- Railway
- Google Cloud Run
- AWS Lambda (with API Gateway)
- Any platform supporting Go applications

## Environment Variables

- `SUPABASE_URL` - Your Supabase project URL
- `SUPABASE_ANON_KEY` - Supabase anonymous key
- `SUPABASE_SERVICE_ROLE_KEY` - Supabase service role key (for backend)
- `PORT` - Server port (default: 8080)
