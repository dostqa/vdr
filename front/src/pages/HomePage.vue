<script setup lang="ts">
import AudioRecorder from '../features/audio/components/AudioRecorder.vue'
import ProcessedAudioBlock from '../features/audio/components/ProcessedAudioBlock.vue'
import { ref } from 'vue'

import type { TranscriptionResponse } from '../features/audio/types/transcription.types'
import { mockResponse } from '../features/audio/services/mockTranscription'
import { uploadAudio } from '../features/audio/services/audio.service'
import TranscriptionBlock from '../features/audio/components/TranscriptionBlock.vue'

import { onMounted } from 'vue'

onMounted(() => {
  transcription.value = mockResponse
})
const waveRef = ref<any>(null)
const processedWaveRef = ref<any>(null)

const onSeek = (time: number) => {
  waveRef.value?.seekToTime(time)
  processedWaveRef.value?.seekToTime(time)
}
const loading = ref(false)

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms))

const pollRequest = async (id: number) => {
  const start = Date.now()
  const maxTime = 300_000
  const interval = 3000

  while (Date.now() - start < maxTime) {
    const res = await fetch(
      `http://85.239.55.254:8080/api/audiofiles/requests/${id}`
    )

    if (res.ok) {
      // если бек уже отдал готовый json
      const text = await res.text()

      if (text) {
        try {
          const data = JSON.parse(text)
          if (data && data.request_id) {
            return data
          }
        } catch {
          // если не JSON или пусто — продолжаем polling
        }
      }
    }

    await sleep(interval)
  }

  throw new Error('Timeout: processing took too long')
}

const transcription = ref<TranscriptionResponse | null>(null)

const processedAudio = ref<string | null>(null)

const onRecorded = async (data: any) => {
  loading.value = true
  transcription.value = null

  try {
    // 1. отправка аудио
    const uploadRes = await uploadAudio(data.blob)
    const id = uploadRes.id

    // 2. polling результата
    const result = await pollRequest(id)

    // 3. обновление UI
    transcription.value = result
    processedAudio.value =
      `http://85.239.55.254:8080/audio/${result.new_file_path}`

  } catch (e) {
    console.error('Processing error:', e)
  } finally {
    loading.value = false
  }
}

</script>

<template>
  <div class="page">
    <div class="recorder-container">
      <img class="logo" src="/src/assets/icons/logo.png"/>
      <p class="subtitle">Анонимизация госовых данных</p>

      <AudioRecorder @recorded="onRecorded" />

      <div v-if="loading" class="loader-overlay">
  <div class="loader-card">
    <div class="spinner"></div>
    <p>Обрабатываем аудио...</p>
    <span class="sub">это может занять несколько секунд</span>
  </div>
</div>

      <ProcessedAudioBlock
        ref="processedWaveRef"
        :audioUrl="processedAudio"
        :objects="transcription?.objects_pdns"
        style="margin-top: 20px;"
      />

      <TranscriptionBlock
  :data="transcription"
  @seek="onSeek"
/>
    </div>
  </div>
</template>