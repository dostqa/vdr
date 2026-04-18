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

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #eef2f6 100%);
  padding: 20px;
}

.recorder-container {
  text-align: center;
}

.title {
  font-size: 32px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 32px;
}
</style>