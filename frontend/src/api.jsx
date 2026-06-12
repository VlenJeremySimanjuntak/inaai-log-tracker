import axios from 'axios';

const API_BASE = 'http://localhost:8081/api';

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});

export const fetchLogs = () => api.get('/logs');
export const createLog = (data) => api.post('/logs', data);
export const updateStatus = (id, status) => api.put(`/logs/${id}/status`, { status });
export const fetchSummary = () => api.get('/ai-summary');