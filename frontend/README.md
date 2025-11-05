# InfluencePower Frontend (TypeScript/React)

Modern React + TypeScript application for the InfluencePower Chicago City Council visualization.

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building  
- **Tailwind CSS** for styling
- **Radix UI** for accessible components
- **Motion** for animations
- **Lucide React** for icons

## Getting Started

### Prerequisites
- Node.js 18+ and npm

### Installation

```bash
# Install dependencies
npm install
```

### Development

```bash
# Start development server (runs on http://localhost:5173)
npm run dev
```

### Build

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

## Environment Variables

Create a `.env` file in the frontend directory:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

For production, update this to your deployed backend URL.

## Project Structure

```
frontend/
├── src/
│   ├── components/      # React components
│   │   ├── ui/          # Reusable UI components
│   │   ├── ChamberView.tsx
│   │   ├── ProfileView.tsx
│   │   ├── HeroSection.tsx
│   │   └── QuickStats.tsx
│   ├── services/        # API service layer
│   │   └── api.ts       # Backend API calls
│   ├── types/           # TypeScript type definitions
│   │   └── index.ts
│   ├── App.tsx          # Main application component
│   ├── main.tsx         # Application entry point
│   └── index.css        # Global styles
├── index.html           # HTML template
├── package.json
├── tsconfig.json
├── vite.config.ts
└── tailwind.config.js
```

## API Integration

The frontend connects to the Go backend API. Make sure the backend is running on `http://localhost:8080` or update the `VITE_API_URL` environment variable.

### API Endpoints Used

- `GET /api/v1/officials` - Fetch all officials
- `GET /api/v1/officials/{id}` - Get official details
- `GET /api/v1/officials/{id}/voting-records` - Get voting history
- `GET /api/v1/officials/{id}/committees` - Get committee memberships
- `GET /api/v1/wards/{ward}/statistics` - Get ward statistics

## Development Notes

- Hot Module Replacement (HMR) enabled for fast development
- TypeScript strict mode enabled for type safety
- ESLint configured for code quality
- Tailwind CSS for utility-first styling
- Responsive design optimized for mobile and desktop

## Deployment

### Vercel (Recommended)
1. Connect your GitHub repository
2. Set `VITE_API_URL` environment variable to your production backend URL
3. Deploy automatically

### Build and Deploy Manually
```bash
npm run build
# Upload dist/ folder to your hosting service
```

## Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build locally
- `npm run lint` - Run ESLint

## License

MIT
