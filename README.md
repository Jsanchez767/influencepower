# InfluencePower - Chicago City Council Visualization# Chicago City Hall - Council Chamber



An interactive web application for visualizing Chicago City Council officials, their voting records, and ward statistics. Features a modern, progressive campaign-style interface inspired by civic engagement design.An interactive visualization of the Chicago City Council Chamber with a seating chart matching the official council layout.



## 🌟 Features## Features



- **Interactive Chamber Visualization**: SVG-based semicircular chamber layout with polaroid-style cards for all 52 officials (Mayor, 50 Aldermen, City Clerk)- **Interactive Seating Chart**: SVG-based semicircular chamber layout with numbered seats (1-46)

- **Comprehensive Profile Pages**: Detailed views with:- **Official Details**: Click any seat to view alderman information, committee assignments, and legislative priorities

  - Legislative Impact Score with metrics- **Party Filtering**: Filter seats by political party affiliation

  - Voting Pattern & Alliance Tracker- **Legislation Tracker**: View active bills and their status

  - Recent Voting Record- **Directory**: Browse all council members with search and filter capabilities

  - Ward-to-Ward Comparison Dashboard- **Responsive Design**: Works on desktop and mobile devices

  - Committee Memberships

- **Party Filtering**: Filter officials by political affiliation (Democrat/Republican)## Layout

- **Hero Section**: Layered Chicago City Hall imagery with branding

- **Quick Stats**: City metrics dashboard (51 officials, 50 wards, 2.7M residents)The chamber features:

- **RESTful API Backend**: Go backend with Supabase PostgreSQL database- **4 arc rows** of seats (46 total aldermen)

- **Mayor's dais** at the bottom center

## 🏗️ Project Structure- **City Clerk** position in the center

- **Cabinet Members** section on the left

```- **Press** section on the right

InfluencePower/

├── index.html              # Main standalone application (production-ready)## Installation

├── vercel.json            # Vercel deployment configuration

├── backend/               # Go backend API```bash

│   ├── main.go           # Server entry point# Install dependencies

│   ├── api/              # API routesnpm install

│   ├── db/               # Database connection & schema

│   ├── handlers/         # Request handlers# Run development server

│   └── models/           # Data modelsnpm run dev

└── frontend/             # React/Vite development version (optional)

    ├── App.tsx# Build for production

    ├── package.jsonnpm run build

    └── components/

```# Preview production build

npm run preview

## 🚀 Quick Start```



### Frontend (Standalone)## Tech Stack



The main application is a self-contained HTML file with React via CDN:- **React 18** with TypeScript

- **Vite** for fast development and building

```bash- **Tailwind CSS** for styling

# Open directly in browser- **Radix UI** for accessible components

open index.html- **Lucide React** for icons

- **Motion** for animations

# Or serve with any static server

python3 -m http.server 8000## Usage

```

The component can be imported and used in any React application:

### Backend (Go + Supabase)

```tsx

1. **Prerequisites**: Go 1.21+, Supabase accountimport ChicagoCityHall from './App'



2. **Setup**:function MyApp() {

```bash  return <ChicagoCityHall />

cd backend}

cp .env.example .env```

# Edit .env with your Supabase credentials

```## Data Structure



3. **Install dependencies**:Officials are defined with:

```bash- Name, position, ward number

go mod download- Political party affiliation

```- Contact information (phone, email)

- Committee assignments

4. **Set up database**:- Policy priorities

   - Go to your Supabase project → SQL Editor- Recent legislation

   - Run the schema from `backend/db/schema.sql`

## Customization

5. **Run the server**:

```bash- Update `officials` array in `App.tsx` to modify council member data

go run main.go- Adjust seat positions by modifying the `rows` configuration in `generateSeatPositions()`

# Server starts on http://localhost:8080- Customize colors by editing the Chicago theme colors in `tailwind.config.js`

```

## Credits

## 📡 API Endpoints

Based on the Chicago City Council structure with 50 aldermen and 1 mayor.

Base URL: `http://localhost:8080/api/v1`

## License

### Officials

- `GET /officials` - Get all officialsMIT

- `GET /officials/{id}` - Get official by ID
- `POST /officials` - Create new official
- `PUT /officials/{id}` - Update official
- `DELETE /officials/{id}` - Delete official
- `GET /officials/party/{party}` - Filter by party
- `GET /officials/ward/{ward}` - Filter by ward

### Voting Records
- `GET /officials/{id}/voting-records` - Get voting records
- `POST /voting-records` - Create voting record

### Ward Statistics
- `GET /wards/{ward}/statistics` - Get ward statistics

### Committees
- `GET /committees` - Get all committees
- `GET /officials/{id}/committees` - Get official committees

### Health Check
- `GET /health` - API health status

## 🎨 Design

- **Typography**: DM Sans (400, 500, 700, 900)
- **Color Scheme**: Progressive campaign aesthetic with red, white, and blue accents
- **Layout**: Semicircular chamber with 3-section distribution
  - Row 1: 4-6-4 seats (radius 240)
  - Row 2: 5-7-5 seats (radius 310)
  - Row 3: 6-7-6 seats (radius 380)
- **Card Sizes**: 
  - Aldermen: 48x64px
  - Mayor: 70x90px
  - City Clerk: 50x65px

## 🗄️ Database Schema

### Tables
- `officials` - City officials with contact info and roles
- `voting_records` - Historical voting data
- `committees` - Committee information
- `official_committees` - Official-committee relationships
- `ward_statistics` - Ward-level metrics (infrastructure, housing, permits)

## 🚢 Deployment

### Vercel (Frontend)
1. Connect your GitHub repository to Vercel
2. Deploy automatically - `vercel.json` is configured
3. Access at your Vercel URL

### Backend Hosting Options
- **Railway**: `railway up`
- **Heroku**: Create app and push
- **Google Cloud Run**: Deploy as container
- **AWS Lambda**: With API Gateway

## 🛠️ Tech Stack

### Frontend
- React 18 (via CDN)
- Babel Standalone (JSX compilation)
- SVG Graphics
- DM Sans Font

### Backend
- Go 1.21+
- Gorilla Mux (routing)
- Supabase/PostgreSQL (database)
- CORS enabled

## 📝 Environment Variables

Create `backend/.env` with:

```env
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_anon_key
SUPABASE_SERVICE_ROLE_KEY=your_service_role_key
PORT=8080
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📄 License

MIT License - feel free to use this project for your own city council visualization!

## 🙏 Credits

Inspired by modern progressive campaign design and civic engagement tools. Built for transparency and accessibility in local government.
