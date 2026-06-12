import axios from 'axios';

// Ganti localhost dengan domain Railway kamu
const API_BASE = 'https://inaai-log-tracker-production.up.railway.app/api';

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});

export const fetchLogs = () => api.get('/logs');
export const createLog = (data) => api.post('/logs', data);
export const updateStatus = (id, status) => api.put(`/logs/${id}/status`, { status });
export const fetchSummary = () => api.get('/summary/latest');