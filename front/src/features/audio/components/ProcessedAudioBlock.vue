<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import WaveSurfer from 'wavesurfer.js'
import RegionsPlugin from 'wavesurfer.js/dist/plugins/regions.esm.js'

import playIcon from '/src/assets/icons/play.png'
import pauseIcon from '/src/assets/icons/pause.png'
import downloadIcon from '/src/assets/icons/download.png'

const props = defineProps<{
  audioUrl: string | null
  objects?: {
    start_time: number
    end_time: number
  }[]
}>()

const container = ref<HTMLDivElement | null>(null)
let wave: WaveSurfer | null = null
let regions: any = null

const isPlaying = ref(false)

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
})

watch(
  () => props.objects,
  () => {
    addRegions()
  }
)

const addRegions = () => {
  if (!regions || !props.objects) return

  regions.clearRegions()

  props.objects.forEach(obj => {
    regions.addRegion({
      start: obj.start_time,
      end: obj.end_time,
      color: 'rgba(0,255,0,0.2)',
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
    
    <p class="anotation">Анонимный</p>
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