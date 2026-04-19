import { ref } from 'vue'

export function useAudio() {
  const audioUrl = ref<string | null>(null)
  const isRecording = ref(false)
  const audioBlob = ref<Blob | null>(null)

  let mediaRecorder: MediaRecorder | null = null
  let chunks: Blob[] = []
  let stream: MediaStream | null = null

  const startRecording = async () => {
    try {
      stream = await navigator.mediaDevices.getUserMedia({ audio: true })

      const mimeType = MediaRecorder.isTypeSupported('audio/webm')
        ? 'audio/webm'
        : 'audio/mp4'

      mediaRecorder = new MediaRecorder(stream, { mimeType })
      chunks = []

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) chunks.push(e.data)
      }

      mediaRecorder.onstop = () => {
        const blob = new Blob(chunks, { type: mimeType })
        if (audioUrl.value) {
          URL.revokeObjectURL(audioUrl.value)
        }
        audioBlob.value = blob
        audioUrl.value = URL.createObjectURL(blob)
        
        if (stream) {
          stream.getTracks().forEach(track => track.stop())
          stream = null
        }
      }

      mediaRecorder.start()
      isRecording.value = true
    } catch (error) {
      console.error('Ошибка доступа к микрофону:', error)
      alert('Не удалось получить доступ к микрофону')
    }
  }

  const stopRecording = () => {
    if (mediaRecorder && mediaRecorder.state === 'recording') {
      mediaRecorder.stop()
      isRecording.value = false
    }
  }

  return {
    audioUrl,
    isRecording,
    startRecording,
    stopRecording,
    audioBlob
  }
}