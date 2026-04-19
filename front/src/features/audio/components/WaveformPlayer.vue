<template>
  <div ref="container" class="wave" />
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import WaveSurfer from 'wavesurfer.js'

const props = defineProps<{
  audioUrl: string | null
}>()

const container = ref<HTMLDivElement | null>(null)
let wave: WaveSurfer | null = null

const isPlaying = ref(false)

const currentTime = ref(0)
const duration = ref(0)

const formatTime = (sec: number) => {
  const m = Math.floor(sec / 60)
  const s = Math.floor(sec % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

const seekToTime = (time: number) => {
  if (!wave) return

  const duration = wave.getDuration()
  wave.seekTo(time / duration)
}

onMounted(() => {
  if (!container.value) return

  wave = WaveSurfer.create({
    container: container.value,
    waveColor: 'var(--wave-color, #cbd5e1)',
    progressColor: 'var(--progress-color, #3b82f6)',
    cursorColor: 'var(--text-primary, #1e293b)',
    barWidth: 3,
    barGap: 3,
    barRadius: 4,
    height: 56,
    normalize: true
  })

  if (props.audioUrl) {
    wave.load(props.audioUrl)
  }

  wave.on('finish', () => {
    isPlaying.value = false
  })

  wave.on('play', () => {
    isPlaying.value = true
  })

  wave.on('pause', () => {
    isPlaying.value = false
  })

  wave.on('ready', () => {
  duration.value = wave?.getDuration() || 0
  currentTime.value = wave?.getCurrentTime() || 0
})

wave.on('audioprocess', () => {
  currentTime.value = wave?.getCurrentTime() || 0
})

wave.on('interaction', () => {
  currentTime.value = wave?.getCurrentTime() || 0
})
})

watch(() => props.audioUrl, (url) => {
  if (wave && url) {
    wave.load(url)
    isPlaying.value = false
  }
})

const playPause = () => {
  if (!wave) return
  wave.playPause()
}

defineExpose({
  playPause,
  isPlaying,
  currentTime,
  duration,
  formatTime,
  seekToTime
})
</script>