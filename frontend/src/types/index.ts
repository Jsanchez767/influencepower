export interface Official {
  id: number;
  name: string;
  ward: number | null;
  party: string;
  role: string;
  contact: string;
  email: string;
  image_url?: string;
  created_at?: string;
  updated_at?: string;
}

export interface VotingRecord {
  id: number;
  official_id: number;
  bill_title: string;
  vote: 'yes' | 'no' | 'abstain';
  vote_date: string;
  description?: string;
  created_at?: string;
}

export interface Committee {
  id: number;
  name: string;
  description?: string;
  created_at?: string;
}

export interface OfficialCommittee {
  id: number;
  official_id: number;
  committee_id: number;
  role: string;
  created_at?: string;
  committee?: Committee;
}

export interface WardStatistic {
  id: number;
  ward: number;
  infrastructure_spending: number;
  response_311_time: number;
  affordable_housing: number;
  permit_processing: number;
  created_at?: string;
  updated_at?: string;
}
