import React, { useState, useEffect } from 'react';
import { fetchOfficials } from './services/api';
import type { Official } from './types';

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

  const filteredOfficials = officials.filter(official => {
    if (filter === 'all') return true;
    return official.party.toLowerCase() === filter;
  });

  if (loading) {
    return <div style={{ padding: '40px', textAlign: 'center' }}><h2>Loading...</h2></div>;
  }

  return (
    <div style={{ fontFamily: 'system-ui, sans-serif', minHeight: '100vh', background: '#f3f4f6' }}>
      <header style={{ background: 'linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%)', color: 'white', padding: '40px 20px', textAlign: 'center' }}>
        <h1 style={{ fontSize: '48px', fontWeight: 900, margin: 0 }}>InfluencePower</h1>
        <p style={{ fontSize: '20px', marginTop: '10px' }}>Chicago City Council</p>
      </header>

      {error && (
        <div style={{ maxWidth: '1200px', margin: '20px auto', padding: '0 20px' }}>
          <div style={{ background: '#fef3c7', border: '2px solid #f59e0b', borderRadius: '8px', padding: '16px' }}>
            ⚠️ {error}
          </div>
        </div>
      )}

      <div style={{ maxWidth: '1200px', margin: '20px auto', padding: '0 20px', display: 'flex', gap: '10px', justifyContent: 'center' }}>
        <button onClick={() => setFilter('all')} style={{ padding: '10px 20px', borderRadius: '8px', border: filter === 'all' ? '2px solid #3b82f6' : '2px solid #e5e7eb', background: filter === 'all' ? '#3b82f6' : 'white', color: filter === 'all' ? 'white' : '#374151', fontWeight: 600, cursor: 'pointer' }}>All</button>
        <button onClick={() => setFilter('democrat')} style={{ padding: '10px 20px', borderRadius: '8px', border: filter === 'democrat' ? '2px solid #3b82f6' : '2px solid #e5e7eb', background: filter === 'democrat' ? '#3b82f6' : 'white', color: filter === 'democrat' ? 'white' : '#374151', fontWeight: 600, cursor: 'pointer' }}>Democrats</button>
        <button onClick={() => setFilter('republican')} style={{ padding: '10px 20px', borderRadius: '8px', border: filter === 'republican' ? '2px solid #3b82f6' : '2px solid #e5e7eb', background: filter === 'republican' ? '#3b82f6' : 'white', color: filter === 'republican' ? 'white' : '#374151', fontWeight: 600, cursor: 'pointer' }}>Republicans</button>
      </div>

      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '0 20px 40px' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
          {filteredOfficials.map(official => (
            <div key={official.id} onClick={() => setSelectedOfficial(official)} style={{ background: 'white', borderRadius: '12px', padding: '20px', boxShadow: '0 2px 8px rgba(0,0,0,0.1)', cursor: 'pointer' }}>
              <div style={{ width: '60px', height: '60px', borderRadius: '50%', background: official.party === 'Democrat' ? '#3b82f6' : '#ef4444', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '24px', fontWeight: 900, margin: '0 auto 16px' }}>
                {official.name.charAt(0)}
              </div>
              <h3 style={{ margin: '0 0 8px 0', fontSize: '18px', textAlign: 'center' }}>{official.name}</h3>
              <p style={{ margin: '0', fontSize: '14px', color: '#666', textAlign: 'center' }}>{official.role}{official.ward ? ` - Ward ${official.ward}` : ''}</p>
            </div>
          ))}
        </div>
      </div>

      {selectedOfficial && (
        <div onClick={() => setSelectedOfficial(null)} style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 1000, padding: '20px' }}>
          <div onClick={(e) => e.stopPropagation()} style={{ background: 'white', borderRadius: '12px', padding: '32px', maxWidth: '500px', width: '100%' }}>
            <h2 style={{ margin: '0 0 10px 0' }}>{selectedOfficial.name}</h2>
            <p style={{ color: '#666', margin: '0 0 20px 0' }}>{selectedOfficial.role}{selectedOfficial.ward && ` - Ward ${selectedOfficial.ward}`}</p>
            <div><strong>Party:</strong> {selectedOfficial.party}</div>
            <div><strong>Phone:</strong> {selectedOfficial.contact}</div>
            <div><strong>Email:</strong> {selectedOfficial.email}</div>
            <button onClick={() => setSelectedOfficial(null)} style={{ marginTop: '20px', padding: '10px 20px', borderRadius: '8px', border: 'none', background: '#3b82f6', color: 'white', fontWeight: 600, cursor: 'pointer', width: '100%' }}>Close</button>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
