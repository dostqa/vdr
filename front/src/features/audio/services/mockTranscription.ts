import type { TranscriptionResponse } from '../types/transcription.types'

export const mockResponse: TranscriptionResponse = {
  request_id: 5501,
  old_file_path: '',
  new_file_path: '',
  original_text: "Мой паспорт серия 4510 номер 123456, телефон 89005553535.",
  anon_text: "Мой паспорт серия [PASSPORT], телефон [PHONE].",
  objects_pdns: [
    {
      text: "4510 123456",
      type: "PASSPORT",
      start_time: 2.15,
      end_time: 4.80
    },
    {
      text: "89005553535",
      type: "PHONE",
      start_time: 6.30,
      end_time: 8.10
    }
  ]
}