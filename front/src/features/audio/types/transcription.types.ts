export interface PDNObject {
  text: string
  type: string
  start_time: number
  end_time: number
}

export interface TranscriptionResponse {
  request_id: number
  old_file_path: string
  new_file_path: string
  original_text: string
  anon_text: string
  objects_pdns: PDNObject[]
}