<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TranscriptionResponse } from '../types/transcription.types'

import eyeIcon from '/src/assets/icons/eye.png'
import eyeOffIcon from '/src/assets/icons/eye-off.png'

const props = defineProps<{
  data: TranscriptionResponse | null
}>()

const showOriginal = ref(false)

/**
 * Подсветка [TYPE] в anon тексте
 */
const highlightedAnonText = computed(() => {
  if (!props.data) return ''

  return props.data.anon_text.replace(
    /\[(.*?)\]/g,
    '<span class="mask">[$1]</span>'
  )
})

const toggle = () => {
  showOriginal.value = !showOriginal.value
}
</script>

<template>
    <p class="anotation">Транскрипция</p>
  <div class="text-block">
    <div class="header">

      <button class="control-btn" @click="toggle" :disabled="!data">
        <img :src="showOriginal ? eyeOffIcon : eyeIcon" width="18" />
      </button>
    </div>

    <div class="text-content">
      <!-- если нет данных -->
      <div v-if="!data" class="empty">
        нет данных
      </div>

      <!-- анонимный -->
      <div
        v-else-if="!showOriginal"
        v-html="highlightedAnonText"
      />

      <!-- оригинал -->
      <div v-else>
        {{ data.original_text }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.text-block {
  margin-top: 20px;
  padding: 16px 20px;
  border-radius: var(--radius-md);
  background: var(--surface-white);
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-sm);
  width: 90%;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
  color: var(--text-secondary);
}

.text-content {
  font-size: 15px;
  color: var(--text-primary);
  line-height: 1.5;
}

:deep(.mask) {
  background: rgba(0,255,0,0.2);
  border-radius: 4px;
  padding: 2px 4px;
}

.empty {
  color: var(--text-secondary);
}
</style>