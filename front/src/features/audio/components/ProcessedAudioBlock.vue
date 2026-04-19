<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import WaveSurfer from 'wavesurfer.js'
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js'

import playIcon from '/src/assets/icons/play.png'
import pauseIcon from '/src/assets/icons/pause.png'
import downloadIcon from '/src/assets/icons/download.png'

const currentTime = ref(0)
const duration = ref(0)

const TYPE_COLORS: Record<string, string> = {
  PASSPORT: 'rgba(255, 99, 132, 0.3)',
  INN: 'rgba(54, 162, 235, 0.3)',
  PHONE: 'rgba(75, 192, 192, 0.3)',
  EMAIL: 'rgba(153, 102, 255, 0.3)',
  ADDRESS: 'rgba(255, 206, 86, 0.3)',
  SNILS: 'rgba(255, 159, 64, 0.3)'
}

const formatTime = (sec: number) => {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

const props = defineProps<{
  audioUrl: string | null
  objects?: {
    start_time: number
    end_time: number
    type: string
  }[]
}>()

const container = ref<HTMLDivElement | null>(null)
let wave: WaveSurfer | null = null
let regions: any = null

const isPlaying = ref(false)

const seekToTime = (time: number) => {
  if (!wave) return

  const duration = wave.getDuration()
  wave.seekTo(time / duration)
}

defineExpose({
  seekToTime
})

onMounted(() => {
  if (!container.value) return

  regions = RegionsPlugin.create()

  wave = WaveSurfer.create({
    container: container.value,
    waveColor: 'var(--wave-color, #cbd5e1)',
    progressColor: 'var(--progress-color, #3b82f6)',
    cursorColor: 'var(--text-primary, #1e293b)',
    barWidth: 3,
    barGap: 3,
    barRadius: 4,
    height: 56,
    normalize: true,
    plugins: [regions]
  })

  if (props.audioUrl) wave.load(props.audioUrl)

  wave.on('ready', () => {
    addRegions()
  })

  wave.on('finish', () => {
    isPlaying.value = false
  })

  wave.on('ready', () => {
  duration.value = wave?.getDuration() || 0
})

wave.on('audioprocess', () => {
  currentTime.value = wave?.getCurrentTime() || 0
})

wave.on('interaction', () => {
  currentTime.value = wave?.getCurrentTime() || 0
})
})

watch(
  () => props.objects,
  () => {
    addRegions()
  }
)

watch(() => props.audioUrl, (url) => {
  if (wave && url) {
    wave.load(url)
    currentTime.value = 0
    duration.value = 0
  }
})

const addRegions = () => {
  if (!regions || !props.objects) return

  regions.clearRegions()

  props.objects.forEach(obj => {
    regions.addRegion({
      start: obj.start_time,
      end: obj.end_time,
      color: TYPE_COLORS[obj.type] || 'rgba(0,255,0,0.2)', // fallback
      drag: false,
      resize: false
    })
  })
}

const togglePlay = () => {
  wave?.playPause()
  isPlaying.value = !isPlaying.value
}

// скачать
const downloadAudio = () => {
  if (!props.audioUrl) return

  const a = document.createElement('a')
  a.href = props.audioUrl
  a.download = 'processed.webm'
  a.click()
}
</script>

<template>
    
    <div class="info-row">
  <p class="anotation">Анонимный</p>
  <span class="anotation">{{ formatTime(currentTime) }} / {{ formatTime(duration) }}</span>
</div>
  <div class="audio-block">
    <!-- play -->
    <button class="control-btn" @click="togglePlay" :disabled="!audioUrl">
      <img :src="isPlaying ? pauseIcon : playIcon" width="15" />
    </button>

    <!-- waveform -->
    <div class="wave-container">
      <div ref="container"/>
    </div>

    <!-- download -->
    <button class="control-btn" @click="downloadAudio" :disabled="!audioUrl">
      <img :src="downloadIcon" width="20" />
    </button>
  </div>
</template>

<style scoped>
.label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}
</style>