import { useState, useEffect } from 'react';
import { fetchOfficials } from './services/api';
import type { Official } from './types';

interface Seat {
  x: number;
  y: number;
  rotation: number;
  official: Official;
}

function App() {
  const [officials, setOfficials] = useState<Official[]>([]);
  const [selectedOfficial, setSelectedOfficial] = useState<Official | null>(null);
  const [filter, setFilter] = useState<'all' | 'democrat' | 'republican'>('all');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadOfficials();
  }, []);

  const loadOfficials = async () => {
    try {
      setLoading(true);
      const data = await fetchOfficials();
      setOfficials(data);
      setError(null);
    } catch (err) {
      setError('Backend not connected. Using mock data.');
      setOfficials(getMockData());
    } finally {
      setLoading(false);
    }
  };

  const getMockData = (): Official[] => {
    return [
      { id: 1, name: 'Brandon Johnson', ward: null, party: 'Democrat', role: 'Mayor', contact: '(312) 744-3300', email: 'mayor@cityofchicago.org' },
      ...Array.from({ length: 50 }, (_, i) => ({
        id: i + 2,
        name: `Alderman ${i + 1}`,
        ward: i + 1,
        party: i % 3 === 0 ? 'Republican' : 'Democrat',
        role: 'Alderman',
        contact: `(312) 744-${3000 + i}`,
        email: `ward${String(i + 1).padStart(2, '0')}@cityofchicago.org`
      })),
      { id: 52, name: 'Anna Valencia', ward: null, party: 'Democrat', role: 'City Clerk', contact: '(312) 744-6861', email: 'clerk@cityofchicago.org' }
    ];
  };

  const generateSeatPositions = (): Seat[] => {
    const seats: Seat[] = [];
    let officialIndex = 1;
    
    const centerX = 450;
    const centerY = 320;
    
    const rows = [
      { radius: 240, leftCount: 4, centerCount: 6, rightCount: 4, leftAisle: 15, rightAisle: 15 },
      { radius: 310, leftCount: 5, centerCount: 7, rightCount: 5, leftAisle: 15, rightAisle: 15 },
      { radius: 380, leftCount: 6, centerCount: 7, rightCount: 6, leftAisle: 15, rightAisle: 15 }
    ];
    
    rows.forEach(row => {
      const totalAngle = 180;
      const leftAngle = (totalAngle - row.leftAisle - row.rightAisle) * 0.28;
      const centerAngle = (totalAngle - row.leftAisle - row.rightAisle) * 0.44;
      const rightAngle = (totalAngle - row.leftAisle - row.rightAisle) * 0.28;
      
      // Left section
      for (let i = 0; i < row.leftCount; i++) {
        const angle = 180 + (i / (row.leftCount - 1)) * leftAngle;
        const rad = (angle * Math.PI) / 180;
        const aldermen = officials.filter(o => o.role === 'Alderman');
        if (aldermen[officialIndex - 1]) {
          seats.push({
            x: centerX + row.radius * Math.cos(rad),
            y: centerY + row.radius * Math.sin(rad),
            rotation: angle - 90,
            official: aldermen[officialIndex - 1]
          });
        }
        officialIndex++;
      }
      
      // Center section
      const centerStartAngle = 180 + leftAngle + row.leftAisle;
      for (let i = 0; i < row.centerCount; i++) {
        const angle = centerStartAngle + (i / (row.centerCount - 1)) * centerAngle;
        const rad = (angle * Math.PI) / 180;
        const aldermen = officials.filter(o => o.role === 'Alderman');
        if (aldermen[officialIndex - 1]) {
          seats.push({
            x: centerX + row.radius * Math.cos(rad),
            y: centerY + row.radius * Math.sin(rad),
            rotation: angle - 90,
            official: aldermen[officialIndex - 1]
          });
        }
        officialIndex++;
      }
      
      // Right section
      const rightStartAngle = centerStartAngle + centerAngle + row.rightAisle;
      for (let i = 0; i < row.rightCount; i++) {
        const angle = rightStartAngle + (i / (row.rightCount - 1)) * rightAngle;
        const rad = (angle * Math.PI) / 180;
        const aldermen = officials.filter(o => o.role === 'Alderman');
        if (aldermen[officialIndex - 1]) {
          seats.push({
            x: centerX + row.radius * Math.cos(rad),
            y: centerY + row.radius * Math.sin(rad),
            rotation: angle - 90,
            official: aldermen[officialIndex - 1]
          });
        }
        officialIndex++;
      }
    });
    
    return seats;
  };

  const mayor = officials.find(o => o.role === 'Mayor');
  const clerk = officials.find(o => o.role === 'City Clerk');
  
  const seats = generateSeatPositions();
  
  const filteredSeats = filter === 'all' 
    ? seats 
    : seats.filter(seat => seat.official.party.toLowerCase() === filter);

  if (loading) {
    return <div style={{ padding: '40px', textAlign: 'center' }}><h2>Loading...</h2></div>;
  }

  return (
    <div style={{ fontFamily: "'DM Sans', system-ui, sans-serif", minHeight: '100vh', background: '#f8f9fa', color: '#1a1a1a' }}>
      {/* Hero Section with City Hall Building */}
      <section style={{ position: 'relative', width: '100%', height: 0, paddingBottom: '56.25%', background: '#E8E5DF' }}>
        <div style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', overflow: 'hidden' }}>
          
          {/* City Hall Building - Base layer */}
          <div style={{ position: 'absolute', width: '100%', height: '100%', zIndex: 10, top: '-25%', transform: 'translateY(25%)' }}>
            <img
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2976ede320a2c6e6e27_city%20hall%20building.png"
              alt="City Hall Building"
              style={{ width: '100%', height: '100%', objectFit: 'contain', objectPosition: 'center center' }}
            />
          </div>

          {/* HALL Text - Below everything */}
          <div style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', zIndex: 0 }}>
            <img
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2975695c32b41340967_text%2C%20hall.png"
              alt="Hall"
              style={{ width: '100%', height: '100%', objectFit: 'contain', objectPosition: 'center center' }}
            />
          </div>

          {/* CITY Text - In front of building */}
          <div style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', zIndex: 20 }}>
            <img
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b297e001e99022c35bc8_text%2C%20City.png"
              alt="City"
              style={{ width: '100%', height: '100%', objectFit: 'contain', objectPosition: 'center center' }}
            />
          </div>

          {/* Chicago Text with Stars */}
          <div style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', zIndex: 10 }}>
            <img
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2977fd3a46a6280b698_Chicago%20w%20stars.png"
              alt="Chicago with stars"
              style={{ width: '100%', height: '100%', objectFit: 'contain', objectPosition: 'center center' }}
            />
          </div>

          {/* Chicago Ribbon */}
          <div style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', zIndex: 30 }}>
            <img
              src="https://cdn.prod.website-files.com/68adde99fad93ee4c3740168/6907b2972e758df565d5f921_Chicago%20ribbon.png"
              alt="Chicago ribbon"
              style={{ width: '100%', height: '100%', objectFit: 'contain', objectPosition: 'center center' }}
            />
          </div>
        </div>
      </section>

      {/* Quick Stats */}
      <section style={{ background: 'white', borderBottom: '4px solid #C8102E', padding: '32px 0' }}>
        <div style={{ maxWidth: '1400px', margin: '0 auto', padding: '0 16px' }}>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '24px' }}>
            <div style={{ textAlign: 'center' }}>
              <div style={{ color: '#C8102E', fontSize: '2.5rem', fontWeight: 900, marginBottom: '8px' }}>
                51
              </div>
              <div style={{ color: '#6b7280', fontSize: '0.875rem', fontWeight: 700, textTransform: 'uppercase' }}>
                Elected Officials
              </div>
            </div>
            <div style={{ textAlign: 'center' }}>
              <div style={{ color: '#C8102E', fontSize: '2.5rem', fontWeight: 900, marginBottom: '8px' }}>
                50
              </div>
              <div style={{ color: '#6b7280', fontSize: '0.875rem', fontWeight: 700, textTransform: 'uppercase' }}>
                City Wards
              </div>
            </div>
            <div style={{ textAlign: 'center' }}>
              <div style={{ color: '#C8102E', fontSize: '2.5rem', fontWeight: 900, marginBottom: '8px' }}>
                2.7M
              </div>
              <div style={{ color: '#6b7280', fontSize: '0.875rem', fontWeight: 700, textTransform: 'uppercase' }}>
                Residents Served
              </div>
            </div>
            <div style={{ textAlign: 'center' }}>
              <div style={{ color: '#C8102E', fontSize: '2.5rem', fontWeight: 900, marginBottom: '8px' }}>
                100%
              </div>
              <div style={{ color: '#6b7280', fontSize: '0.875rem', fontWeight: 700, textTransform: 'uppercase' }}>
                Transparency
              </div>
            </div>
          </div>
        </div>
      </section>

      <div style={{ maxWidth: '1400px', margin: '0 auto', padding: '40px 20px' }}>
        {/* Header / Title Section */}
        <div style={{ textAlign: 'center', marginBottom: '60px' }}>
          <h1 style={{ fontSize: '2.5rem', fontWeight: 900, letterSpacing: '-0.02em', color: '#1a1a1a', marginBottom: '16px' }}>
            The City Council of the City of <span style={{ color: '#2563eb' }}>Chicago</span>
          </h1>
          <p style={{ fontSize: '1.125rem', color: '#6b7280', fontWeight: 500 }}>
            50 Aldermen representing Chicago's diverse communities
          </p>
        </div>

        {error && (
          <div style={{ maxWidth: '1200px', margin: '20px auto', padding: '0 20px' }}>
            <div style={{ background: '#fef3c7', border: '2px solid #f59e0b', borderRadius: '8px', padding: '16px' }}>
              ⚠️ {error}
            </div>
          </div>
        )}

        {/* Filter Buttons */}
        <div style={{ display: 'flex', justifyContent: 'center', gap: '12px', marginBottom: '40px', flexWrap: 'wrap' }}>
          <button 
            onClick={() => setFilter('all')} 
            style={{ 
              padding: '12px 24px', 
              fontSize: '0.9375rem', 
              fontWeight: 700, 
              border: filter === 'all' ? '2px solid #2563eb' : '2px solid #e5e7eb', 
              background: filter === 'all' ? '#2563eb' : 'white', 
              color: filter === 'all' ? 'white' : '#6b7280', 
              borderRadius: '12px', 
              cursor: 'pointer',
              fontFamily: "'DM Sans', sans-serif",
              letterSpacing: '-0.01em',
              boxShadow: filter === 'all' ? '0 4px 12px rgba(37,99,235,0.3)' : 'none'
            }}
          >
            All Members
          </button>
          <button 
            onClick={() => setFilter('democrat')} 
            style={{ 
              padding: '12px 24px', 
              fontSize: '0.9375rem', 
              fontWeight: 700, 
              border: filter === 'democrat' ? '2px solid #2563eb' : '2px solid #e5e7eb', 
              background: filter === 'democrat' ? '#2563eb' : 'white', 
              color: filter === 'democrat' ? 'white' : '#6b7280', 
              borderRadius: '12px', 
              cursor: 'pointer',
              fontFamily: "'DM Sans', sans-serif",
              letterSpacing: '-0.01em',
              boxShadow: filter === 'democrat' ? '0 4px 12px rgba(37,99,235,0.3)' : 'none'
            }}
          >
            Democrats
          </button>
          <button 
            onClick={() => setFilter('republican')} 
            style={{ 
              padding: '12px 24px', 
              fontSize: '0.9375rem', 
              fontWeight: 700, 
              border: filter === 'republican' ? '2px solid #2563eb' : '2px solid #e5e7eb', 
              background: filter === 'republican' ? '#2563eb' : 'white', 
              color: filter === 'republican' ? 'white' : '#6b7280', 
              borderRadius: '12px', 
              cursor: 'pointer',
              fontFamily: "'DM Sans', sans-serif",
              letterSpacing: '-0.01em',
              boxShadow: filter === 'republican' ? '0 4px 12px rgba(37,99,235,0.3)' : 'none'
            }}
          >
            Republicans
          </button>
        </div>

        {/* Chamber Container */}
        <div style={{ background: 'white', borderRadius: '24px', padding: '30px 20px', boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -1px rgba(0,0,0,0.06)', marginBottom: '40px' }}>
          <div style={{ width: '100%', margin: '0 auto', display: 'flex', justifyContent: 'center' }}>
            <svg viewBox="0 -150 900 600" xmlns="http://www.w3.org/2000/svg" style={{ maxWidth: '100%', height: 'auto' }}>
              {filteredSeats.map((seat, idx) => {
                const cardWidth = 48;
                const cardHeight = 64;
                
                return (
                  <g
                    key={idx}
                    transform={`translate(${seat.x}, ${seat.y}) rotate(${seat.rotation})`}
                    onClick={() => { setSelectedOfficial(seat.official); }}
                    style={{ cursor: 'pointer' }}
                  >
                    <rect
                      x={-cardWidth/2}
                      y={-cardHeight/2}
                      width={cardWidth}
                      height={cardHeight}
                      rx="3"
                      fill="white"
                      stroke="#e5e7eb"
                      strokeWidth="3"
                      style={{ filter: 'drop-shadow(0 1px 3px rgba(0,0,0,0.12))' }}
                    />
                    
                    <image
                      href={`https://i.pravatar.cc/150?img=${seat.official.id + 10}`}
                      x={-cardWidth/2 + 3}
                      y={-cardHeight/2 + 3}
                      width={cardWidth - 6}
                      height={cardWidth - 6}
                      style={{ clipPath: 'inset(0 round 4px)' }}
                    />
                    
                    <text
                      x="0"
                      y={cardHeight/2 - 3}
                      textAnchor="middle"
                      style={{ fontFamily: "'DM Sans', sans-serif", fontSize: '6px', fontWeight: 600, fill: '#6b7280' }}
                    >
                      Ward {seat.official.ward}
                    </text>
                  </g>
                );
              })}

              {/* Mayor Card */}
              {mayor && (
                <g
                  transform="translate(450, 380)"
                  onClick={() => { setSelectedOfficial(mayor); }}
                  style={{ cursor: 'pointer' }}
                >
                  <rect 
                    x="-35" 
                    y="-45" 
                    width="70" 
                    height="90" 
                    rx="5"
                    fill="#3b82f6"
                    stroke="#2563eb"
                    strokeWidth="3"
                    style={{ filter: 'drop-shadow(0 1px 3px rgba(0,0,0,0.12))' }}
                  />
                  <image
                    href={`https://i.pravatar.cc/150?img=33`}
                    x="-32"
                    y="-42"
                    width="64"
                    height="64"
                    style={{ clipPath: 'inset(0 round 4px)' }}
                  />
                  <text 
                    x="0" 
                    y="35" 
                    textAnchor="middle" 
                    style={{ fontFamily: "'DM Sans', sans-serif", fontSize: '9px', fontWeight: 900, fill: 'white' }}
                  >
                    Mayor
                  </text>
                  <text 
                    x="0" 
                    y="43" 
                    textAnchor="middle" 
                    style={{ fontFamily: "'DM Sans', sans-serif", fontSize: '8px', fontWeight: 600, fill: 'white' }}
                  >
                    {mayor.name.split(' ')[0]}
                  </text>
                </g>
              )}

              {/* City Clerk Card */}
              {clerk && (
                <g
                  transform="translate(450, 270)"
                  onClick={() => { setSelectedOfficial(clerk); }}
                  style={{ cursor: 'pointer' }}
                >
                  <rect 
                    x="-25" 
                    y="-30" 
                    width="50" 
                    height="65" 
                    rx="4"
                    fill="#10b981"
                    stroke="#059669"
                    strokeWidth="3"
                    style={{ filter: 'drop-shadow(0 1px 3px rgba(0,0,0,0.12))' }}
                  />
                  <image
                    href={`https://i.pravatar.cc/150?img=45`}
                    x="-22"
                    y="-27"
                    width="44"
                    height="44"
                    style={{ clipPath: 'inset(0 round 4px)' }}
                  />
                  <text 
                    x="0" 
                    y="28" 
                    textAnchor="middle" 
                    style={{ fontFamily: "'DM Sans', sans-serif", fontSize: '7px', fontWeight: 700, fill: 'white' }}
                  >
                    City Clerk
                  </text>
                  <text 
                    x="0" 
                    y="34" 
                    textAnchor="middle" 
                    style={{ fontFamily: "'DM Sans', sans-serif", fontSize: '6px', fontWeight: 600, fill: 'white' }}
                  >
                    {clerk.name.split(' ')[0]}
                  </text>
                </g>
              )}
            </svg>
          </div>
        </div>
      </div>

      {/* Modal */}
      {selectedOfficial && (
        <div 
          onClick={() => setSelectedOfficial(null)} 
          style={{ 
            position: 'fixed', 
            inset: 0, 
            background: 'rgba(0,0,0,0.7)', 
            backdropFilter: 'blur(4px)', 
            zIndex: 50, 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center', 
            padding: '20px' 
          }}
        >
          <div 
            onClick={(e) => e.stopPropagation()} 
            style={{ 
              background: 'white', 
              borderRadius: '20px', 
              maxWidth: '600px', 
              width: '100%', 
              maxHeight: '90vh', 
              overflowY: 'auto', 
              boxShadow: '0 20px 60px rgba(0,0,0,0.3)', 
              position: 'relative' 
            }}
          >
            {/* Modal Header */}
            <div style={{ padding: '32px 32px 24px', borderBottom: '1px solid #e5e7eb' }}>
              <button 
                onClick={() => setSelectedOfficial(null)} 
                style={{ 
                  position: 'absolute', 
                  top: '20px', 
                  right: '20px', 
                  background: '#f3f4f6', 
                  border: 'none', 
                  width: '36px', 
                  height: '36px', 
                  borderRadius: '50%', 
                  cursor: 'pointer', 
                  display: 'flex', 
                  alignItems: 'center', 
                  justifyContent: 'center', 
                  fontSize: '20px', 
                  color: '#6b7280' 
                }}
              >
                ×
              </button>
              
              <img 
                src={`https://i.pravatar.cc/150?img=${selectedOfficial.id + 10}`}
                alt={selectedOfficial.name}
                style={{ width: '120px', height: '120px', borderRadius: '16px', objectFit: 'cover', marginBottom: '20px' }}
              />
              
              <h2 style={{ fontSize: '1.75rem', fontWeight: 900, color: '#1a1a1a', marginBottom: '8px', letterSpacing: '-0.02em' }}>
                {selectedOfficial.name}
              </h2>
              <p style={{ fontSize: '1.125rem', color: '#6b7280', fontWeight: 600, marginBottom: '24px' }}>
                {selectedOfficial.role}{selectedOfficial.ward && ` - Ward ${selectedOfficial.ward}`}
              </p>
            </div>

            {/* Modal Body */}
            <div style={{ padding: '32px' }}>
              <div style={{ display: 'grid', gap: '16px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '12px', background: '#f9fafb', borderRadius: '12px' }}>
                  <span style={{ fontSize: '0.875rem', fontWeight: 700, color: '#6b7280', textTransform: 'uppercase', letterSpacing: '0.05em', minWidth: '100px' }}>Party</span>
                  <span 
                    style={{ 
                      display: 'inline-block', 
                      padding: '6px 16px', 
                      borderRadius: '8px', 
                      fontSize: '0.875rem', 
                      fontWeight: 700, 
                      letterSpacing: '0.02em',
                      background: selectedOfficial.party === 'Democrat' ? '#dbeafe' : '#fee2e2',
                      color: selectedOfficial.party === 'Democrat' ? '#1e40af' : '#991b1b'
                    }}
                  >
                    {selectedOfficial.party}
                  </span>
                </div>
                
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '12px', background: '#f9fafb', borderRadius: '12px' }}>
                  <span style={{ fontSize: '0.875rem', fontWeight: 700, color: '#6b7280', textTransform: 'uppercase', letterSpacing: '0.05em', minWidth: '100px' }}>Phone</span>
                  <span style={{ fontSize: '1rem', fontWeight: 600, color: '#1a1a1a' }}>{selectedOfficial.contact}</span>
                </div>
                
                <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '12px', background: '#f9fafb', borderRadius: '12px' }}>
                  <span style={{ fontSize: '0.875rem', fontWeight: 700, color: '#6b7280', textTransform: 'uppercase', letterSpacing: '0.05em', minWidth: '100px' }}>Email</span>
                  <span style={{ fontSize: '1rem', fontWeight: 600, color: '#1a1a1a' }}>{selectedOfficial.email}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
