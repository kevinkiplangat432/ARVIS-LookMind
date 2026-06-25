import axios from 'axios'
import type { Request, Anomaly } from '../types'

const http = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

export async function fetchRequests(limit = 50): Promise<Request[]> {
  const res = await http.get<Request[]>('/requests', { params: { limit } })
  return res.data
}

export async function fetchAnomalies(limit = 50): Promise<Anomaly[]> {
  const res = await http.get<Anomaly[]>('/anomalies', { params: { limit } })
  return res.data
}