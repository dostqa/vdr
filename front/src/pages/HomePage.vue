<script setup lang="ts">
import AudioRecorder from '../features/audio/components/AudioRecorder.vue'
import ProcessedAudioBlock from '../features/audio/components/ProcessedAudioBlock.vue'
import { ref } from 'vue'

const processedAudio = ref<string | null>(null)

const onRecorded = (url: string) => {
  setTimeout(() => {
    processedAudio.value = url
  }, 1000)
}
</script>

<template>
  <div class="page">
    <div class="recorder-container">
      <img class="logo" src="/src/assets/icons/logo.png"/>
      <p class="subtitle">Анонимизация госовых данных</p>

      <AudioRecorder @recorded="onRecorded" />

      <ProcessedAudioBlock
        :audioUrl="processedAudio"
        style="margin-top: 20px;"
      />
    </div>
  </div>
</template>