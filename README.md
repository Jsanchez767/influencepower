# Chicago City Hall - Council Chamber

An interactive visualization of the Chicago City Council Chamber with a seating chart matching the official council layout.

## Features

- **Interactive Seating Chart**: SVG-based semicircular chamber layout with numbered seats (1-46)
- **Official Details**: Click any seat to view alderman information, committee assignments, and legislative priorities
- **Party Filtering**: Filter seats by political party affiliation
- **Legislation Tracker**: View active bills and their status
- **Directory**: Browse all council members with search and filter capabilities
- **Responsive Design**: Works on desktop and mobile devices

## Layout

The chamber features:
- **4 arc rows** of seats (46 total aldermen)
- **Mayor's dais** at the bottom center
- **City Clerk** position in the center
- **Cabinet Members** section on the left
- **Press** section on the right

## Installation

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Tailwind CSS** for styling
- **Radix UI** for accessible components
- **Lucide React** for icons
- **Motion** for animations

## Usage

The component can be imported and used in any React application:

```tsx
import ChicagoCityHall from './App'

function MyApp() {
  return <ChicagoCityHall />
}
```

## Data Structure

Officials are defined with:
- Name, position, ward number
- Political party affiliation
- Contact information (phone, email)
- Committee assignments
- Policy priorities
- Recent legislation

## Customization

- Update `officials` array in `App.tsx` to modify council member data
- Adjust seat positions by modifying the `rows` configuration in `generateSeatPositions()`
- Customize colors by editing the Chicago theme colors in `tailwind.config.js`

## Credits

Based on the Chicago City Council structure with 50 aldermen and 1 mayor.

## License

MIT
