export interface Request {
  id: string
  model: string
  prompt_tokens: number
  completion_tokens: number
  latency_ms: number
  status_code: number
  created_at: string
}

export interface Anomaly {
  id: string
  request_id: string
  rule: string
  detail: string
  created_at: string
}

export interface Stats {
  total_requests: number
  total_anomalies: number
  avg_latency_ms: number
}