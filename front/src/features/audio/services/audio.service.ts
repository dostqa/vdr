import type { TranscriptionResponse } from '../types/transcription.types'

export async function uploadAudio(file: Blob): Promise<{ id: number }> {
  const formData = new FormData()
  formData.append('file', file)

  const response = await fetch('audiofiles', {
    method: 'POST',
    body: formData
  })

  if (!response.ok) {
    throw new Error('Ошибка загрузки аудио')
  }

  return response.json()
}