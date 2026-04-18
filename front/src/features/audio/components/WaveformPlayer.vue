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

onMounted(() => {
  if (!container.value) return

  wave = WaveSurfer.create({
    container: container.value,
    waveColor: '#cbd5e1',
    progressColor: '#3b82f6',
    cursorColor: '#1e293b',
    barWidth: 3,
    barGap: 3,
    barRadius: 4,
    height: 56,
    responsive: true,
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
  isPlaying
})
</script>

<style scoped>
.wave {
  width: 100%;
  cursor: pointer;
}

.wave :deep(wave) {
  border-radius: 40px;
}
</style>