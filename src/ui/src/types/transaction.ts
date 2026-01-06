/**
 * Transaction types matching the backend API structure
 */

export interface Transaction {
  id: string
  external_id: string
  amount: number
  currency: string
  origin: string
  destination: string
  type: string
  status: 'pending' | 'approved' | 'rejected' | 'in_review'
  processing_time_ms: number
  matched_rules: string[]
  metadata: Record<string, any>
  created_at: string
  processed_at: string | null
}
