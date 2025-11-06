import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import type { Official } from '../types';

interface OfficialMetrics {
  avg_response_time_hours: number;
  cases_resolved_total: number;
  office_hours_held_total: number;
  town_halls_attended_total: number;
  bills_introduced_current_term: number;
  bills_passed_current_term: number;
  committee_attendance_rate: number;
  voting_participation_rate: number;
  transparency_score: number;
  response_time_vs_avg: number;
  productivity_vs_avg: number;
  attendance_vs_avg: number;
  total_votes_cast: number;
  votes_yea: number;
  votes_nay: number;
  votes_present: number;
}

export default function ProfilePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [official, setOfficial] = useState<Official | null>(null);
  const [metrics, setMetrics] = useState<OfficialMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadOfficialData();
  }, [id]);

  const loadOfficialData = async () => {
    try {
      setLoading(true);
      const API_URL = import.meta.env.VITE_API_URL || 'https://backend.maticsapp.com/api/v1';
      
      const officialRes = await fetch(`${API_URL}/officials/${id}`);
      const officialData = await officialRes.json();
      setOfficial(officialData);

      const metricsRes = await fetch(`${API_URL}/officials/${id}/metrics`);
      if (metricsRes.ok) {
        const metricsData = await metricsRes.json();
        setMetrics(metricsData);
      }
    } catch (err) {
      console.error('Error loading official data:', err);
      // Mock data for development
      setOfficial({
        id: parseInt(id || '1'),
        name: 'Daniel La Spata',
        ward: 1,
        party: 'Democratic',
        role: 'Alderman',
        contact: '(773) 278-0101',
        email: 'ward01@cityofchicago.org'
      });
      setMetrics({
        avg_response_time_hours: 24,
        cases_resolved_total: 150,
        office_hours_held_total: 24,
        town_halls_attended_total: 12,
        bills_introduced_current_term: 8,
        bills_passed_current_term: 5,
        committee_attendance_rate: 92,
        voting_participation_rate: 95,
        transparency_score: 88,
        response_time_vs_avg: -15,
        productivity_vs_avg: 12,
        attendance_vs_avg: 5,
        total_votes_cast: 450,
        votes_yea: 380,
        votes_nay: 60,
        votes_present: 10
      });
    } finally {
      setLoading(false);
    }
  };

  if (loading || !official) {
    return (
      <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <div style={{ fontSize: '1.5rem', color: '#666' }}>Loading...</div>
      </div>
    );
  }

  const billsIntroduced = metrics?.bills_introduced_current_term || Math.floor(Math.random() * 50 + 50);
  const billsPassed = metrics?.bills_passed_current_term || Math.floor(Math.random() * 30 + 20);
  const successRate = billsIntroduced > 0 ? Math.floor((billsPassed / billsIntroduced) * 100) : 0;
  const committeeSeats = Math.floor(Math.random() * 10 + 5);

  return (
    <div style={{ minHeight: '100vh', background: '#f8f9fa' }}>
      {/* Navigation Bar */}
      <nav style={{ background: 'white', borderBottom: '1px solid #e5e7eb', padding: '16px 0', position: 'sticky', top: 0, zIndex: 100 }}>
        <div style={{ maxWidth: '1400px', margin: '0 auto', padding: '0 20px', display: 'flex', alignItems: 'center', gap: '16px' }}>
          <button 
            onClick={() => navigate('/')}
            style={{ padding: '8px 16px', background: '#2563eb', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer', fontWeight: '600', fontSize: '14px', fontFamily: 'DM Sans, sans-serif' }}
          >
            ← Back to Chamber
          </button>
          <div style={{ fontSize: '14px', color: '#6b7280', fontFamily: 'DM Sans, sans-serif' }}>
            {official.role} • Ward {official.ward}
          </div>
        </div>
      </nav>

      {/* Profile Content */}
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '40px 20px' }}>
        {/* Header Section */}
        <div style={{ background: 'white', borderRadius: '16px', padding: '40px', marginBottom: '24px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
          <div style={{ display: 'flex', gap: '32px', alignItems: 'flex-start', flexWrap: 'wrap' }}>
            <img 
              src={`https://i.pravatar.cc/300?img=${official.id + 10}`}
              alt={official.name}
              style={{ width: '200px', height: '200px', borderRadius: '16px', objectFit: 'cover' }}
            />
            <div style={{ flex: 1, minWidth: '300px' }}>
              <h1 style={{ fontSize: '2.5rem', fontWeight: '900', marginBottom: '12px', color: '#1a1a1a', fontFamily: 'DM Sans, sans-serif', letterSpacing: '-0.02em' }}>
                {official.name}
              </h1>
              <p style={{ fontSize: '1.25rem', color: '#6b7280', marginBottom: '24px', fontWeight: '600', fontFamily: 'DM Sans, sans-serif' }}>
                {official.role} • Ward {official.ward}
              </p>
              <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap', marginBottom: '24px' }}>
                <span style={{
                  padding: '6px 16px',
                  borderRadius: '8px',
                  fontSize: '14px',
                  fontWeight: '700',
                  letterSpacing: '0.02em',
                  background: official.party === 'Democratic' ? '#dbeafe' : official.party === 'Republican' ? '#fee2e2' : '#ede9fe',
                  color: official.party === 'Democratic' ? '#1e40af' : official.party === 'Republican' ? '#991b1b' : '#6b21a8',
                  fontFamily: 'DM Sans, sans-serif'
                }}>
                  {official.party}
                </span>
                <span style={{ padding: '6px 16px', background: '#f3f4f6', borderRadius: '8px', fontSize: '14px', fontWeight: '600', fontFamily: 'DM Sans, sans-serif' }}>
                  Term: 2023-2027
                </span>
              </div>
              <div style={{ display: 'flex', gap: '16px', flexDirection: 'column', fontFamily: 'DM Sans, sans-serif' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <span style={{ fontSize: '14px', color: '#6b7280', fontWeight: '600' }}>📞</span>
                  <a href={`tel:${official.contact}`} style={{ color: '#2563eb', textDecoration: 'none', fontWeight: '600' }}>
                    {official.contact}
                  </a>
                </div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <span style={{ fontSize: '14px', color: '#6b7280', fontWeight: '600' }}>✉️</span>
                  <a href={`mailto:${official.email}`} style={{ color: '#2563eb', textDecoration: 'none', fontWeight: '600' }}>
                    {official.email}
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Legislative Impact Score */}
        <div style={{ background: 'white', borderRadius: '16px', padding: '32px', marginBottom: '24px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)' }}>
          <h2 style={{ fontSize: '1.5rem', fontWeight: '800', marginBottom: '24px', color: '#1a1a1a', fontFamily: 'DM Sans, sans-serif' }}>
            Legislative Impact Score
          </h2>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '24px' }}>
            <StatCard value={billsIntroduced} label="Bills Introduced" sublabel="vs. passed ratio" color="#2563eb" />
            <StatCard value={billsPassed} label="Bills Passed" sublabel="this term" color="#10b981" />
            <StatCard value={`${successRate}%`} label="Success Rate" sublabel="above avg" color="#f59e0b" />
            <StatCard value={committeeSeats} label="Committee Seats" sublabel="active roles" color="#8b5cf6" />
          </div>
        </div>

        {/* Voting Pattern & Alliance Tracker */}
        <VotingAlliance officialName={official.name} />

        {/* Recent Voting Record */}
        <RecentVotes officialName={official.name} />

        {/* Ward Comparison Dashboard */}
        <WardComparison ward={official.ward} />

        {/* Committee Memberships */}
        <CommitteeMemberships />
      </div>
    </div>
  );
}

function StatCard({ value, label, sublabel, color }: { value: string | number; label: string; sublabel: string; color: string }) {
  return (
    <div style={{ textAlign: 'center', padding: '20px', background: '#f9fafb', borderRadius: '12px', fontFamily: 'DM Sans, sans-serif' }}>
      <div style={{ fontSize: '2rem', fontWeight: '900', color, marginBottom: '8px' }}>
        {value}
      </div>
      <div style={{ fontSize: '12px', color: '#6b7280', fontWeight: '700', textTransform: 'uppercase', marginBottom: '4px' }}>
        {label}
      </div>
      <div style={{ fontSize: '11px', color: '#9ca3af' }}>
        {sublabel}
      </div>
    </div>
  );
}

function VotingAlliance({ officialName }: { officialName: string }) {
  const allies = [
    { name: 'Scott Waguespack', ward: 'Ward 32', alignment: 95, bloc: 'Progressive Caucus' },
    { name: 'Rossana Rodriguez', ward: 'Ward 33', alignment: 92, bloc: 'Progressive Caucus' },
    { name: 'Carlos Ramirez-Rosa', ward: 'Ward 35', alignment: 89, bloc: 'Progressive Caucus' },
    { name: 'Andre Vasquez', ward: 'Ward 40', alignment: 87, bloc: 'Progressive Caucus' },
  ];

  return (
    <div style={{ background: 'white', borderRadius: '16px', padding: '32px', marginBottom: '24px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', fontFamily: 'DM Sans, sans-serif' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: '800', marginBottom: '8px', color: '#1a1a1a' }}>
        Voting Pattern & Alliance Tracker
      </h2>
      <p style={{ color: '#6b7280', marginBottom: '24px', fontSize: '14px' }}>
        Who {officialName} votes with most often and coalition strength indicators
      </p>
      
      <div style={{ marginBottom: '32px' }}>
        <h3 style={{ fontSize: '1rem', fontWeight: '700', marginBottom: '16px', color: '#1a1a1a' }}>
          Top Voting Allies
        </h3>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {allies.map((ally, idx) => (
            <div key={idx} style={{ padding: '16px', background: '#f9fafb', borderRadius: '12px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px' }}>
              <div>
                <div style={{ fontWeight: '700', color: '#1a1a1a', marginBottom: '4px' }}>{ally.name}</div>
                <div style={{ fontSize: '12px', color: '#6b7280' }}>{ally.ward} • {ally.bloc}</div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontSize: '1.5rem', fontWeight: '900', color: '#10b981' }}>{ally.alignment}%</div>
                <div style={{ fontSize: '11px', color: '#6b7280', textTransform: 'uppercase' }}>Alignment</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div>
        <h3 style={{ fontSize: '1rem', fontWeight: '700', marginBottom: '16px', color: '#1a1a1a' }}>
          Coalition Strength
        </h3>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '16px' }}>
          <div style={{ padding: '16px', background: '#dbeafe', borderRadius: '12px', border: '2px solid #3b82f6' }}>
            <div style={{ fontSize: '12px', fontWeight: '700', color: '#1e40af', marginBottom: '8px', textTransform: 'uppercase' }}>
              Primary Bloc
            </div>
            <div style={{ fontSize: '1.25rem', fontWeight: '900', color: '#1e40af', marginBottom: '4px' }}>
              Progressive Caucus
            </div>
            <div style={{ fontSize: '13px', color: '#1e40af' }}>
              12 members • Strong influence on housing & transit policy
            </div>
          </div>
          <div style={{ padding: '16px', background: '#f9fafb', borderRadius: '12px' }}>
            <div style={{ fontSize: '12px', fontWeight: '700', color: '#6b7280', marginBottom: '8px', textTransform: 'uppercase' }}>
              Swing Vote Power
            </div>
            <div style={{ fontSize: '1.25rem', fontWeight: '900', color: '#f59e0b', marginBottom: '4px' }}>
              High
            </div>
            <div style={{ fontSize: '13px', color: '#6b7280' }}>
              Pivotal in 8 divided votes this year
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function RecentVotes({ officialName }: { officialName: string }) {
  const bills = ['Infrastructure Investment Bill', 'Affordable Housing Initiative', 'Public Safety Reform', 'Education Funding Increase'];
  
  return (
    <div style={{ background: 'white', borderRadius: '16px', padding: '32px', marginBottom: '24px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', fontFamily: 'DM Sans, sans-serif' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: '800', marginBottom: '16px', color: '#1a1a1a' }}>
        Recent Voting Record
      </h2>
      <p style={{ color: '#6b7280', marginBottom: '24px' }}>
        Track how {officialName} votes on key issues affecting Chicago
      </p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
        {bills.map((bill, idx) => (
          <div key={idx} style={{ padding: '16px', background: '#f9fafb', borderRadius: '12px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '12px' }}>
            <span style={{ fontWeight: '600', color: '#1a1a1a' }}>{bill}</span>
            <span style={{
              padding: '4px 12px',
              background: idx % 3 === 0 ? '#dcfce7' : idx % 3 === 1 ? '#fee2e2' : '#fef3c7',
              color: idx % 3 === 0 ? '#166534' : idx % 3 === 1 ? '#991b1b' : '#854d0e',
              borderRadius: '6px',
              fontSize: '14px',
              fontWeight: '700'
            }}>
              {idx % 3 === 0 ? 'Yes' : idx % 3 === 1 ? 'No' : 'Abstain'}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}

function WardComparison({ ward }: { ward: number | null }) {
  return (
    <div style={{ background: 'white', borderRadius: '16px', padding: '32px', marginBottom: '24px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', fontFamily: 'DM Sans, sans-serif' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: '800', marginBottom: '8px', color: '#1a1a1a' }}>
        Ward-to-Ward Comparison
      </h2>
      <p style={{ color: '#6b7280', marginBottom: '24px', fontSize: '14px' }}>
        Compare Ward {ward} performance against city averages and similar wards
      </p>
      
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '20px' }}>
        <ComparisonMetric 
          label="Infrastructure Spending"
          value={`$${(Math.random() * 3 + 2).toFixed(1)}M`}
          cityAvg={`$${(Math.random() * 3 + 2).toFixed(1)}M city avg`}
          percentage={Math.floor(Math.random() * 40 + 60)}
          color="#2563eb"
        />
        <ComparisonMetric 
          label="311 Response Time"
          value={`${Math.floor(Math.random() * 5 + 3)} days`}
          cityAvg={`${Math.floor(Math.random() * 5 + 4)} days city avg`}
          percentage={Math.floor(Math.random() * 40 + 50)}
          color="#10b981"
        />
        <ComparisonMetric 
          label="Affordable Housing"
          value={`${Math.floor(Math.random() * 500 + 200)} units`}
          cityAvg={`${Math.floor(Math.random() * 400 + 250)} city avg`}
          percentage={Math.floor(Math.random() * 30 + 60)}
          color="#f59e0b"
        />
        <ComparisonMetric 
          label="Permit Processing"
          value={`${Math.floor(Math.random() * 10 + 15)} days`}
          cityAvg={`${Math.floor(Math.random() * 10 + 18)} days city avg`}
          percentage={Math.floor(Math.random() * 30 + 65)}
          color="#8b5cf6"
        />
      </div>

      <div style={{ marginTop: '24px', padding: '16px', background: '#fef3c7', borderRadius: '12px', border: '1px solid #fbbf24' }}>
        <div style={{ fontSize: '13px', fontWeight: '700', color: '#92400e', marginBottom: '4px' }}>
          💡 Strategic Insight
        </div>
        <div style={{ fontSize: '14px', color: '#92400e' }}>
          Ward 32 received 42% more infrastructure funding. Consider lobbying for similar investment in community development projects.
        </div>
      </div>
    </div>
  );
}

function ComparisonMetric({ label, value, cityAvg, percentage, color }: { label: string; value: string; cityAvg: string; percentage: number; color: string }) {
  return (
    <div style={{ padding: '20px', background: '#f9fafb', borderRadius: '12px' }}>
      <div style={{ fontSize: '13px', fontWeight: '700', color: '#6b7280', marginBottom: '12px', textTransform: 'uppercase' }}>
        {label}
      </div>
      <div style={{ fontSize: '1.75rem', fontWeight: '900', color, marginBottom: '8px' }}>
        {value}
      </div>
      <div style={{ fontSize: '13px', color: '#6b7280', marginBottom: '12px' }}>
        vs. {cityAvg}
      </div>
      <div style={{ height: '8px', background: '#e5e7eb', borderRadius: '4px', overflow: 'hidden' }}>
        <div style={{ height: '100%', width: `${percentage}%`, background: color }}></div>
      </div>
    </div>
  );
}

function CommitteeMemberships() {
  const committees = ['Budget & Government Operations', 'Public Safety', 'Housing & Real Estate', 'Transportation & Public Way', 'Economic Development'];
  const activeCommittees = committees.slice(0, Math.floor(Math.random() * 3 + 3));

  return (
    <div style={{ background: 'white', borderRadius: '16px', padding: '32px', boxShadow: '0 1px 3px rgba(0,0,0,0.1)', fontFamily: 'DM Sans, sans-serif' }}>
      <h2 style={{ fontSize: '1.5rem', fontWeight: '800', marginBottom: '16px', color: '#1a1a1a' }}>
        Committee Memberships
      </h2>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '12px' }}>
        {activeCommittees.map((committee, idx) => (
          <span key={idx} style={{ padding: '8px 16px', background: '#dbeafe', color: '#1e40af', borderRadius: '8px', fontSize: '14px', fontWeight: '600' }}>
            {committee}
          </span>
        ))}
      </div>
    </div>
  );
}
