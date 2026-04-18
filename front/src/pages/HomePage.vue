<script setup lang="ts">
import AudioRecorder from '../features/audio/components/AudioRecorder.vue'
import ProcessedAudioBlock from '../features/audio/components/ProcessedAudioBlock.vue'
import { ref } from 'vue'

import TranscriptionBlock from '../features/audio/components/TranscriptionBlock.vue'
import { mockResponse } from '../features/audio/services/mockTranscription'

const transcription = ref(null)

const onRecorded = (url: string) => {
  setTimeout(() => {
    processedAudio.value = url

    // имитация ответа бэка
    transcription.value = mockResponse
  }, 1000)
}

const processedAudio = ref<string | null>(null)

</script>

<template>
  <div class="page">
    <div class="recorder-container">
      <img class="logo" src="/src/assets/icons/logo.png"/>
      <p class="subtitle">Анонимизация госовых данных</p>

      <AudioRecorder @recorded="onRecorded" />

      <ProcessedAudioBlock
        :audioUrl="processedAudio"
        :objects="transcription?.objects_pdns"
        style="margin-top: 20px;"
      />

      <TranscriptionBlock :data="transcription" />
    </div>
  </div>
</template>