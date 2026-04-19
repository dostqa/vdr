<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TranscriptionResponse } from '../types/transcription.types'

import eyeIcon from '/src/assets/icons/eye.png'
import eyeOffIcon from '/src/assets/icons/eye-off.png'

const hoveredIndex = ref<number | null>(null)

const props = defineProps<{
  data: TranscriptionResponse | null
}>()

const showOriginal = ref(false)

const emit = defineEmits(['seek'])

const handleClick = (e: MouseEvent) => {
  const target = e.target as HTMLElement

  if (target.classList.contains('mask')) {
    const index = target.dataset.index
    if (!props.data || index === undefined) return

    const obj = props.data.objects_pdns[Number(index)]
  }
}

const handleMouseOver = (e: MouseEvent) => {
  const el = e.target as HTMLElement
  if (el.classList.contains('mask')) {
    const index = el.dataset.index
    if (index !== undefined) {
      hoveredIndex.value = Number(index)
    }
  }
}

const handleMouseLeave = () => {
  hoveredIndex.value = null
}

/**
 * Подсветка [TYPE] в anon тексте
 */
const highlightedAnonText = computed(() => {
  if (!props.data) return ''

  let text = props.data.anon_text

  props.data.objects_pdns.forEach((obj, index) => {
    const placeholder = `[${obj.type}]`

    const displayText =
      hoveredIndex.value === index
        ? obj.text
        : placeholder

    text = text.replace(
      placeholder,
      `<span 
        class="mask mask-${obj.type.toLowerCase()}"
        data-index="${index}"
      >${displayText}</span>`
    )
  })

  return text
})

const toggle = () => {
  showOriginal.value = !showOriginal.value
}
</script>

<template>

        <div class="info-row"><p class="anotation">Транскрипция</p></div>
    
  <div class="text-block">
    <div class="header">

      <button class="control-btn" @click="toggle" :disabled="!data">
        <img :src="showOriginal ? eyeOffIcon : eyeIcon" width="18" />
      </button>
    </div>

    <div class="text-content">
  <div v-if="!data" class="empty">
    нет данных
  </div>

  <div
  v-else-if="!showOriginal"
  v-html="highlightedAnonText"
  @mouseover="handleMouseOver"
  @mouseleave="handleMouseLeave"
  @click="handleClick"
/>

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
  border-radius: 4px;
  padding: 2px 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

/* разные цвета */

:deep(.mask-passport) {
  background: rgba(255, 99, 132, 0.2);
}

:deep(.mask-inn) {
  background: rgba(54, 162, 235, 0.2);
}

:deep(.mask-phone) {
  background: rgba(75, 192, 192, 0.2);
}

:deep(.mask-email) {
  background: rgba(153, 102, 255, 0.2);
}

:deep(.mask-address) {
  background: rgba(255, 206, 86, 0.2);
}

:deep(.mask-snils) {
  background: rgba(255, 159, 64, 0.2);
}

:deep(.mask:hover) {
  transform: scale(1.05);
  filter: brightness(0.9);
}

.empty {
  color: var(--text-secondary);
}
</style>