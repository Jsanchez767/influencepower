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

interface WardMetrics {
  population: number;
  median_income: number;
  avg_service_request_resolution_days: number;
  total_service_requests: number;
  service_requests_resolved: number;
}

export default function ProfilePage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [official, setOfficial] = useState<Official | null>(null);
  const [metrics, setMetrics] = useState<OfficialMetrics | null>(null);
  const [wardMetrics, setWardMetrics] = useState<WardMetrics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadOfficialData();
  }, [id]);

  const loadOfficialData = async () => {
    try {
      setLoading(true);
      const API_URL = import.meta.env.VITE_API_URL || 'https://backend.maticsapp.com/api/v1';
      
      // Fetch official data
      const officialRes = await fetch(`${API_URL}/officials/${id}`);
      const officialData = await officialRes.json();
      setOfficial(officialData);

      // Fetch metrics
      const metricsRes = await fetch(`${API_URL}/officials/${id}/metrics`);
      if (metricsRes.ok) {
        const metricsData = await metricsRes.json();
        setMetrics(metricsData);
      }

      // Fetch ward metrics if applicable
      if (officialData.ward) {
        const wardRes = await fetch(`${API_URL}/wards/${officialData.ward}/metrics`);
        if (wardRes.ok) {
          const wardData = await wardRes.json();
          setWardMetrics(wardData);
        }
      }
    } catch (err) {
      console.error('Error loading official data:', err);
      // Use mock data
      setOfficial({
        id: parseInt(id || '1'),
        name: 'Loading...',
        ward: null,
        party: 'Democratic',
        role: 'Alderman',
        contact: 'N/A',
        email: 'N/A'
      });
      // Mock metrics
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

  return (
    <div style={{ minHeight: '100vh', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
      {/* Header */}
      <div style={{ 
        background: 'rgba(255,255,255,0.1)', 
        backdropFilter: 'blur(10px)',
        borderBottom: '1px solid rgba(255,255,255,0.2)',
        padding: '1rem 2rem'
      }}>
        <button
          onClick={() => navigate('/')}
          style={{
            background: 'rgba(255,255,255,0.2)',
            border: 'none',
            color: 'white',
            padding: '0.5rem 1rem',
            borderRadius: '8px',
            cursor: 'pointer',
            fontSize: '1rem',
            fontWeight: 600
          }}
        >
          ← Back to Chamber
        </button>
      </div>

      {/* Profile Content */}
      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '2rem' }}>
        {/* Official Header */}
        <div style={{
          background: 'white',
          borderRadius: '16px',
          padding: '2rem',
          marginBottom: '2rem',
          boxShadow: '0 10px 40px rgba(0,0,0,0.1)'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '2rem' }}>
            <img
              src={`https://i.pravatar.cc/150?img=${official.id + 10}`}
              alt={official.name}
              style={{
                width: '150px',
                height: '150px',
                borderRadius: '50%',
                objectFit: 'cover',
                border: '4px solid #667eea'
              }}
            />
            <div style={{ flex: 1 }}>
              <h1 style={{ margin: '0 0 0.5rem 0', fontSize: '2.5rem', color: '#1a1a1a' }}>
                {official.name}
              </h1>
              <p style={{ margin: '0 0 1rem 0', fontSize: '1.25rem', color: '#666' }}>
                {official.role}{official.ward && ` - Ward ${official.ward}`}
              </p>
              <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
                <span style={{
                  padding: '0.5rem 1rem',
                  borderRadius: '20px',
                  fontSize: '0.875rem',
                  fontWeight: 600,
                  background: official.party === 'Democratic' ? '#dbeafe' : '#fee2e2',
                  color: official.party === 'Democratic' ? '#1e40af' : '#991b1b'
                }}>
                  {official.party}
                </span>
                {metrics && (
                  <span style={{
                    padding: '0.5rem 1rem',
                    borderRadius: '20px',
                    fontSize: '0.875rem',
                    fontWeight: 600,
                    background: metrics.transparency_score >= 80 ? '#d1fae5' : '#fef3c7',
                    color: metrics.transparency_score >= 80 ? '#065f46' : '#92400e'
                  }}>
                    Transparency: {metrics.transparency_score.toFixed(0)}/100
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Metrics Grid */}
        {metrics && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem', marginBottom: '2rem' }}>
            {/* Constituent Service Metrics */}
            <MetricCard
              title="Constituent Service"
              icon="📞"
              metrics={[
                { label: 'Avg Response Time', value: `${metrics.avg_response_time_hours}h`, comparison: metrics.response_time_vs_avg },
                { label: 'Cases Resolved', value: metrics.cases_resolved_total },
                { label: 'Office Hours Held', value: metrics.office_hours_held_total },
                { label: 'Town Halls Attended', value: metrics.town_halls_attended_total }
              ]}
            />

            {/* Legislative Productivity */}
            <MetricCard
              title="Legislative Productivity"
              icon="📄"
              metrics={[
                { label: 'Bills Introduced (Current Term)', value: metrics.bills_introduced_current_term, comparison: metrics.productivity_vs_avg },
                { label: 'Bills Passed (Current Term)', value: metrics.bills_passed_current_term },
                { label: 'Committee Attendance', value: `${metrics.committee_attendance_rate}%`, comparison: metrics.attendance_vs_avg },
                { label: 'Voting Participation', value: `${metrics.voting_participation_rate}%` }
              ]}
            />

            {/* Voting Record */}
            <MetricCard
              title="Voting Record"
              icon="🗳️"
              metrics={[
                { label: 'Total Votes Cast', value: metrics.total_votes_cast },
                { label: 'Votes Yea', value: `${metrics.votes_yea} (${((metrics.votes_yea / metrics.total_votes_cast) * 100).toFixed(0)}%)` },
                { label: 'Votes Nay', value: `${metrics.votes_nay} (${((metrics.votes_nay / metrics.total_votes_cast) * 100).toFixed(0)}%)` },
                { label: 'Present/Abstain', value: metrics.votes_present }
              ]}
            />
          </div>
        )}

        {/* Ward Comparison */}
        {official.ward && wardMetrics && (
          <div style={{
            background: 'white',
            borderRadius: '16px',
            padding: '2rem',
            marginBottom: '2rem',
            boxShadow: '0 10px 40px rgba(0,0,0,0.1)'
          }}>
            <h2 style={{ margin: '0 0 1.5rem 0', fontSize: '1.5rem', color: '#1a1a1a' }}>
              🗺️ Ward {official.ward} Demographics & Service
            </h2>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1.5rem' }}>
              <StatBox label="Population" value={wardMetrics.population.toLocaleString()} />
              <StatBox label="Median Income" value={`$${wardMetrics.median_income.toLocaleString()}`} />
              <StatBox label="Avg Service Resolution" value={`${wardMetrics.avg_service_request_resolution_days} days`} />
              <StatBox label="Service Requests" value={`${wardMetrics.service_requests_resolved}/${wardMetrics.total_service_requests}`} />
            </div>
          </div>
        )}

        {/* Performance vs City Average */}
        {metrics && (
          <div style={{
            background: 'white',
            borderRadius: '16px',
            padding: '2rem',
            boxShadow: '0 10px 40px rgba(0,0,0,0.1)'
          }}>
            <h2 style={{ margin: '0 0 1.5rem 0', fontSize: '1.5rem', color: '#1a1a1a' }}>
              📊 Performance vs City Average
            </h2>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))', gap: '1rem' }}>
              <ComparisonBar
                label="Response Time"
                percentage={metrics.response_time_vs_avg}
                inversed={true}
              />
              <ComparisonBar
                label="Legislative Productivity"
                percentage={metrics.productivity_vs_avg}
              />
              <ComparisonBar
                label="Committee Attendance"
                percentage={metrics.attendance_vs_avg}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface MetricCardProps {
  title: string;
  icon: string;
  metrics: Array<{ label: string; value: string | number; comparison?: number }>;
}

function MetricCard({ title, icon, metrics }: MetricCardProps) {
  return (
    <div style={{
      background: 'white',
      borderRadius: '16px',
      padding: '1.5rem',
      boxShadow: '0 10px 40px rgba(0,0,0,0.1)'
    }}>
      <h3 style={{ margin: '0 0 1rem 0', fontSize: '1.25rem', color: '#1a1a1a', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
        <span style={{ fontSize: '1.5rem' }}>{icon}</span>
        {title}
      </h3>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        {metrics.map((metric, idx) => (
          <div key={idx}>
            <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.25rem' }}>
              {metric.label}
            </div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
              <span style={{ fontSize: '1.5rem', fontWeight: 700, color: '#1a1a1a' }}>
                {metric.value}
              </span>
              {metric.comparison !== undefined && (
                <span style={{
                  fontSize: '0.875rem',
                  fontWeight: 600,
                  color: metric.comparison > 0 ? '#059669' : '#dc2626'
                }}>
                  {metric.comparison > 0 ? '+' : ''}{metric.comparison.toFixed(0)}% vs avg
                </span>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

interface StatBoxProps {
  label: string;
  value: string;
}

function StatBox({ label, value }: StatBoxProps) {
  return (
    <div>
      <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.5rem' }}>
        {label}
      </div>
      <div style={{ fontSize: '1.75rem', fontWeight: 700, color: '#1a1a1a' }}>
        {value}
      </div>
    </div>
  );
}

interface ComparisonBarProps {
  label: string;
  percentage: number;
  inversed?: boolean;
}

function ComparisonBar({ label, percentage, inversed = false }: ComparisonBarProps) {
  const isPositive = inversed ? percentage < 0 : percentage > 0;
  const absPercentage = Math.abs(percentage);
  const barWidth = Math.min(absPercentage, 100);

  return (
    <div>
      <div style={{ fontSize: '0.875rem', color: '#666', marginBottom: '0.5rem' }}>
        {label}
      </div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
        <div style={{ flex: 1, height: '32px', background: '#f3f4f6', borderRadius: '8px', overflow: 'hidden', position: 'relative' }}>
          <div style={{
            position: 'absolute',
            left: 0,
            top: 0,
            height: '100%',
            width: `${barWidth}%`,
            background: isPositive ? 'linear-gradient(90deg, #10b981, #059669)' : 'linear-gradient(90deg, #ef4444, #dc2626)',
            transition: 'width 0.3s ease'
          }} />
          <div style={{
            position: 'absolute',
            left: 0,
            top: 0,
            width: '100%',
            height: '100%',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '0.875rem',
            fontWeight: 700,
            color: barWidth > 30 ? 'white' : '#1a1a1a'
          }}>
            {percentage > 0 ? '+' : ''}{percentage.toFixed(0)}%
          </div>
        </div>
      </div>
    </div>
  );
}
