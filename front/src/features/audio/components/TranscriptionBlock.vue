<script setup lang="ts">
import { ref, computed } from 'vue'
import type { TranscriptionResponse } from '../types/transcription.types'

import eyeIcon from '/src/assets/icons/eye.png'
import eyeOffIcon from '/src/assets/icons/eye-off.png'

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

/**
 * Подсветка [TYPE] в anon тексте
 */
const highlightedAnonText = computed(() => {
  if (!props.data) return ''

  let text = props.data.anon_text

  props.data.objects_pdns.forEach((obj, index) => {
    const placeholder = `[${obj.type}]`

    text = text.replace(
      placeholder,
      `<span 
        class="mask"
        data-index="${index}"
        title="${obj.text}"
      >${placeholder}</span>`
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
    @click="handleClick"
  ></div>

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