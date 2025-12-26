/**
 * Transaction types matching the backend API structure
 */

export interface Transaction {
  id: string
  external_id: string
  amount: number
  currency: string
  from_account: string
  to_account: string
  type: string
  status: 'pending' | 'approved' | 'rejected' | 'review'
  risk_score: number
  risk_level: 'low' | 'medium' | 'high'
  processing_time_ms: number
  matched_rules: string[]
  metadata: Record<string, any>
  created_at: string
  processed_at: string | null
}
