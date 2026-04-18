<template>
  <div class="audio-block">
    <!-- play кнопка -->
    <button 
      class="control-btn play-btn" 
      @click="togglePlay" 
      :disabled="!audioUrl"
    >
      <svg v-if="!waveRef?.isPlaying" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polygon points="5 3 19 12 5 21 5 3" fill="currentColor"/>
      </svg>
      <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="6" y="4" width="4" height="16" fill="currentColor"/>
        <rect x="14" y="4" width="4" height="16" fill="currentColor"/>
      </svg>
    </button>

    <!-- waveform -->
    <div class="wave-container">
      <WaveformPlayer
        ref="waveRef"
        :audioUrl="audioUrl"
      />
    </div>

    <!-- mic кнопка -->
    <button
      class="control-btn mic-btn"
      :class="{ active: isRecording }"
      @click="toggleRecord"
    >
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M12 1C10.8954 1 10 1.89543 10 3V12C10 13.1046 10.8954 14 12 14C13.1046 14 14 13.1046 14 12V3C14 1.89543 13.1046 1 12 1Z" fill="currentColor"/>
        <path d="M19 10V12C19 15.866 15.866 19 12 19C8.13401 19 5 15.866 5 12V10" stroke="currentColor" fill="none"/>
        <line x1="12" y1="19" x2="12" y2="23" stroke="currentColor"/>
        <line x1="9" y1="23" x2="15" y2="23" stroke="currentColor"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAudio } from '../composables/useAudio'
import WaveformPlayer from './WaveformPlayer.vue'

const {
  audioUrl,
  isRecording,
  startRecording,
  stopRecording
} = useAudio()

const waveRef = ref<any>(null)

const toggleRecord = async () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    await startRecording()
  }
}

const togglePlay = () => {
  waveRef.value?.playPause()
}
</script>

<style scoped>
.audio-block {
  display: flex;
  align-items: center;
  gap: 16px;
  background: #ffffff;
  padding: 16px 24px;
  border-radius: 48px;
  width: 100%;
  max-width: 680px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
}

.audio-block:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  border-color: #cbd5e1;
}

/* общие стили для кнопок */
.control-btn {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

/* play кнопка */
.play-btn {
  background: #f1f5f9;
  color: #1e293b;
}

.play-btn:hover:not(:disabled) {
  background: #e2e8f0;
  transform: scale(1.05);
}

.play-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.play-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* mic кнопка */
.mic-btn {
  background: #f1f5f9;
  color: #1e293b;
}

.mic-btn:hover {
  background: #e2e8f0;
  transform: scale(1.05);
}

.mic-btn:active {
  transform: scale(0.98);
}

.mic-btn.active {
  background: #ef4444;
  color: white;
  animation: pulse 1.5s infinite;
}

/* waveform контейнер */
.wave-container {
  flex: 1;
  min-width: 0;
}

/* анимация записи */
@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.4);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(239, 68, 68, 0);
  }
}
</style>