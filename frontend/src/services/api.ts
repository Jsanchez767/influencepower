import { Official, VotingRecord, Committee, OfficialCommittee, WardStatistic } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Officials
export const fetchOfficials = async (): Promise<Official[]> => {
  const response = await fetch(`${API_BASE_URL}/officials`);
  if (!response.ok) throw new Error('Failed to fetch officials');
  return response.json();
};

export const fetchOfficialById = async (id: number): Promise<Official> => {
  const response = await fetch(`${API_BASE_URL}/officials/${id}`);
  if (!response.ok) throw new Error('Failed to fetch official');
  return response.json();
};

export const fetchOfficialsByParty = async (party: string): Promise<Official[]> => {
  const response = await fetch(`${API_BASE_URL}/officials/party/${party}`);
  if (!response.ok) throw new Error('Failed to fetch officials by party');
  return response.json();
};

export const fetchOfficialsByWard = async (ward: number): Promise<Official[]> => {
  const response = await fetch(`${API_BASE_URL}/officials/ward/${ward}`);
  if (!response.ok) throw new Error('Failed to fetch officials by ward');
  return response.json();
};

// Voting Records
export const fetchVotingRecords = async (officialId: number): Promise<VotingRecord[]> => {
  const response = await fetch(`${API_BASE_URL}/officials/${officialId}/voting-records`);
  if (!response.ok) throw new Error('Failed to fetch voting records');
  return response.json();
};

// Ward Statistics
export const fetchWardStatistics = async (ward: number): Promise<WardStatistic> => {
  const response = await fetch(`${API_BASE_URL}/wards/${ward}/statistics`);
  if (!response.ok) throw new Error('Failed to fetch ward statistics');
  return response.json();
};

// Committees
export const fetchCommittees = async (): Promise<Committee[]> => {
  const response = await fetch(`${API_BASE_URL}/committees`);
  if (!response.ok) throw new Error('Failed to fetch committees');
  return response.json();
};

export const fetchOfficialCommittees = async (officialId: number): Promise<OfficialCommittee[]> => {
  const response = await fetch(`${API_BASE_URL}/officials/${officialId}/committees`);
  if (!response.ok) throw new Error('Failed to fetch official committees');
  return response.json();
};

// Health check
export const checkHealth = async (): Promise<{ status: string; message: string }> => {
  const response = await fetch(`${API_BASE_URL}/health`);
  if (!response.ok) throw new Error('API health check failed');
  return response.json();
};
