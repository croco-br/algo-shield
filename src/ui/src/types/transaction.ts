/**
 * Transaction types matching the backend API structure
 * All transaction fields come from the schema and are stored in metadata
 */

export interface Transaction {
  id: string
  schema_id?: string
  schema_name?: string
  status: 'pending' | 'approved' | 'rejected' | 'in_review'
  processing_time_ms: number
  matched_rules: string[]
  metadata: Record<string, any>
  created_at: string
  processed_at: string | null
}
