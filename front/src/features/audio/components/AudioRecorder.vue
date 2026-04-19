<template>
  <div>

  <div class="info-row">
  <p class="anotation">Исходный</p>
  <span class="anotation">{{ timeDisplay }}</span>
</div>

    <div class="audio-block">
      <button 
        class="control-btn play-btn" 
        @click="togglePlay" 
        :disabled="!audioUrl"
      >
        <img 
          v-if="!waveRef?.isPlaying" 
          :src="playIcon" 
          alt="Play" 
          width="15" 
          height="15"
        />
        <img 
          v-else 
          :src="pauseIcon" 
          alt="Pause" 
          width="20" 
          height="20"
        />
      </button>

      <div class="wave-container">
        <WaveformPlayer
          ref="waveRef"
          :audioUrl="audioUrl"
        />
      </div>

      <button
        class="control-btn mic-btn"
        :class="{ active: isRecording }"
        @click="toggleRecord"
      >
        <img 
          :src="isRecording ? micActiveIcon : micIcon" 
          alt="Microphone" 
          width="20" 
          height="20"
        />
      </button>

      <button class="control-btn" @click="handleUpload">
        <img :src="uploadIcon" width="20" />
      </button>

      <button class="control-btn" @click="downloadAudio" :disabled="!audioUrl">
        <img :src="downloadIcon" width="20"/>
      </button>

      <input
        type="file"
        accept="audio/*"
        ref="fileInput"
        style="display: none"
        @change="onFileChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAudio } from '../composables/useAudio'
import WaveformPlayer from './WaveformPlayer.vue'

import playIcon from '/src/assets/icons/play.png'
import pauseIcon from '/src/assets/icons/pause.png'
import micIcon from '/src/assets/icons/mic.png'
import micActiveIcon from '/src/assets/icons/mic-active.png'
import uploadIcon from '/src/assets/icons/upload.png'
import downloadIcon from '/src/assets/icons/download.png'

import { watch } from 'vue'

import { computed } from 'vue'

const timeDisplay = computed(() => {
  if (!waveRef.value) return '0:00 / 0:00'

  const current = waveRef.value.currentTime || 0
  const total = waveRef.value.duration || 0

  const format = (s: number) => {
    const m = Math.floor(s / 60)
    const sec = Math.floor(s % 60)
    return `${m}:${sec.toString().padStart(2, '0')}`
  }

  return `${format(current)} / ${format(total)}`
})

const emit = defineEmits(['recorded'])

const fileInput = ref<HTMLInputElement | null>(null)

const handleUpload = () => {
  fileInput.value?.click()
}

const onFileChange = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) {
    const url = URL.createObjectURL(file)
    audioUrl.value = url

    emit('recorded', {
      url,
      blob: file
    })
  }
}



const downloadAudio = () => {
  if (!audioUrl.value) return
  const a = document.createElement('a')
  a.href = audioUrl.value
  a.download = 'recording.webm'
  a.click()
}

const {
  audioUrl,
  audioBlob,
  isRecording,
  startRecording,
  stopRecording
} = useAudio()

watch(audioUrl, (val) => {
  if (val && audioBlob.value) {
    emit('recorded', {
      url: val,
      blob: audioBlob.value
    })
  }
})

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