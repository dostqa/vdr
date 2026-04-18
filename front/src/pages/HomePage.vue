<template>
  <div class="page">
    <div class="recorder-container">
      <h1 class="title">Voice Recorder</h1>
      <p class="subtitle">Запиши и прослушай свой голос</p>
      <AudioRecorder />
    </div>
  </div>
</template>

<script setup lang="ts">
import AudioRecorder from '../features/audio/components/AudioRecorder.vue'
</script>