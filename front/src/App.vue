<template>
  <div class="container">
    <!-- Хедер -->
    <header class="header">
      <h1 class="title">🎙️ Voice Recorder</h1>
      <p class="subtitle">Запиши свой голос в высоком качестве</p>
    </header>

    <!-- Переключатель стилей -->
    <div class="style-selector">
      <h3 class="selector-title">🎨 Выбери стиль визуализации</h3>
      <div class="style-buttons">
        <button
          v-for="style in visualStyles"
          :key="style.id"
          @click="currentStyle = style.id"
          class="style-btn"
          :class="{ active: currentStyle === style.id }"
        >
          <span class="style-icon">{{ style.icon }}</span>
          <span>{{ style.name }}</span>
        </button>
      </div>
    </div>

    <!-- Основная зона записи -->
    <div class="recorder-area">
      <div class="visualizer-container">
        <canvas ref="canvasRef" class="visualizer-canvas"></canvas>
        <div class="visualizer-overlay" v-if="!isRecording && !audioUrl">
          <span>🎤 Нажми на кнопку и говори</span>
        </div>
      </div>

      <button @click="toggleRecording" class="record-btn" :class="{ recording: isRecording }">
        <div class="btn-content">
          <span class="mic-icon">{{ isRecording ? '⏹️' : '🎙️' }}</span>
          <span>{{ isRecording ? 'Остановить запись' : 'Начать запись' }}</span>
        </div>
        <div class="recording-pulse" v-if="isRecording"></div>
      </button>

      <div class="timer" v-if="isRecording || recordingDuration > 0">
        <div class="timer-circle">
          <span class="timer-text">{{ formatTime(recordingDuration) }}</span>
        </div>
      </div>
    </div>

    <!-- Кастомный плеер -->
    <div class="playback-area" v-if="audioUrl">
      <div class="playback-card">
        <h3>Фронтендер постарался нормально</h3>

        <div class="waveform-container">
          <canvas ref="waveformCanvasRef" class="waveform-canvas" @click="seekAudio"></canvas>

          <div class="progress-overlay">
            <div class="progress-bar" :style="{ width: playbackProgress + '%' }"></div>
          </div>

          <div class="playhead" :style="{ left: playbackProgress + '%' }">
            <div class="playhead-dot"></div>
          </div>
        </div>

        <div class="custom-controls">
          <button @click="togglePlayback" class="play-btn">
            <span class="play-icon">{{ isPlaying ? '⏸' : '▶' }}</span>
            <span>{{ isPlaying ? 'Пауза' : 'Воспроизвести' }}</span>
          </button>

          <div class="time-info">
            <span>{{ formatTime(currentTime) }}</span>
            <span>/</span>
            <span>{{ formatTime(duration) }}</span>
          </div>

          <button @click="downloadAudio" class="download-btn">Сохранить</button>

          <button @click="resetRecording" class="reset-btn">Новая запись</button>
        </div>
      </div>
    </div>

    <div class="tips" v-if="!audioUrl && !isRecording">
      <div class="tip-item">✨ Нажми на микрофон</div>
      <div class="tip-item">🎤 Говори в микрофон</div>
      <div class="tip-item">⏹️ Нажми ещё раз для остановки</div>
      <div class="tip-item">🎨 Выбери стиль визуализации</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'

const visualStyles = [
  { id: 'bars', name: '🎵 Бары', icon: '📊' },
  { id: 'wave', name: '🌊 Волна', icon: '〰️' },
  { id: 'circular', name: '🔄 Круг', icon: '⚪' },
  { id: 'particles', name: '✨ Частицы', icon: '💫' },
  { id: 'fire', name: '🔥 Огонь', icon: '🔥' },
]

const currentStyle = ref<string>('bars')

// ========== Реактивные переменные ==========
const isRecording = ref<boolean>(false)
const isPlaying = ref<boolean>(false)
const audioUrl = ref<string>('')
const audioBlob = ref<Blob | null>(null)
const recordingDuration = ref<number>(0)
const currentTime = ref<number>(0)
const duration = ref<number>(0)
const playbackProgress = ref<number>(0)

// Refs
const canvasRef = ref<HTMLCanvasElement | null>(null)
const waveformCanvasRef = ref<HTMLCanvasElement | null>(null)

// Аудио элементы (используем стандартный Audio для надёжности)
let audioElement: HTMLAudioElement | null = null
let mediaRecorder: MediaRecorder | null = null
let audioChunks: Blob[] = []
let durationInterval: number | null = null

// Для визуализации записи
let recordingContext: AudioContext | null = null
let analyserNode: AnalyserNode | null = null
let sourceNode: MediaStreamAudioSourceNode | null = null
let animationId: number | null = null
let mediaStream: MediaStream | null = null

// Данные волны
let waveformData: number[] = []
let progressUpdateInterval: number | null = null

// ========== Вспомогательные функции ==========
function formatTime(seconds: number): string {
  if (isNaN(seconds) || seconds < 0) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

// ========== Генерация визуализации дорожки ==========
async function generateWaveformFromBlob(blob: Blob) {
  return new Promise((resolve) => {
    const fileReader = new FileReader()
    fileReader.onloadend = async () => {
      const arrayBuffer = fileReader.result as ArrayBuffer
      const audioContext = new AudioContext()

      try {
        const audioBuffer = await audioContext.decodeAudioData(arrayBuffer)
        duration.value = audioBuffer.duration

        // Получаем данные канала
        const channelData = audioBuffer.getChannelData(0)
        const samples = 200
        waveformData = []
        const blockSize = Math.floor(channelData.length / samples)

        for (let i = 0; i < samples; i++) {
          let sum = 0
          for (let j = 0; j < blockSize; j++) {
            const index = i * blockSize + j
            if (index < channelData.length) {
              sum += Math.abs(channelData[index])
            }
          }
          let average = sum / blockSize
          waveformData.push(Math.min(1, average * 5))
        }

        drawWaveform()
        resolve(true)
      } catch (err) {
        console.error('Ошибка декодирования:', err)
        // Fallback: создаём псевдо-волну
        duration.value = recordingDuration.value
        waveformData = Array(200)
          .fill(0)
          .map(() => Math.random() * 0.5 + 0.2)
        drawWaveform()
        resolve(false)
      } finally {
        audioContext.close()
      }
    }
    fileReader.readAsArrayBuffer(blob)
  })
}

function drawWaveform() {
  if (!waveformCanvasRef.value || waveformData.length === 0) return

  const canvas = waveformCanvasRef.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  canvas.width = canvas.clientWidth
  canvas.height = canvas.clientHeight

  const width = canvas.width
  const height = canvas.height
  const barWidth = width / waveformData.length

  ctx.fillStyle = '#0a0e27'
  ctx.fillRect(0, 0, width, height)

  for (let i = 0; i < waveformData.length; i++) {
    const barHeight = waveformData[i] * height
    const x = i * barWidth
    const y = (height - barHeight) / 2

    const gradient = ctx.createLinearGradient(x, y, x, y + barHeight)
    gradient.addColorStop(0, '#667eea')
    gradient.addColorStop(0.5, '#764ba2')
    gradient.addColorStop(1, '#f093fb')

    ctx.fillStyle = gradient
    ctx.fillRect(x, y, barWidth - 1, barHeight)
  }
}

// ========== Визуализация записи ==========
function drawRecordingVisualization(dataArray: Uint8Array) {
  if (!canvasRef.value) return

  const canvas = canvasRef.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  canvas.width = canvas.clientWidth
  canvas.height = canvas.clientHeight

  ctx.fillStyle = '#0a0e27'
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  const bufferLength = dataArray.length
  const barWidth = (canvas.width / bufferLength) * 2

  switch (currentStyle.value) {
    case 'bars':
      for (let i = 0; i < bufferLength; i++) {
        const barHeight = (dataArray[i] / 255) * canvas.height
        ctx.fillStyle = `hsl(${200 + dataArray[i] / 2}, 70%, 60%)`
        ctx.fillRect(i * barWidth, canvas.height - barHeight, barWidth - 1, barHeight)
      }
      break

    case 'wave':
      ctx.beginPath()
      ctx.strokeStyle = '#667eea'
      ctx.lineWidth = 3
      for (let i = 0; i < bufferLength; i++) {
        const y = (dataArray[i] / 255) * canvas.height
        const x = (i / bufferLength) * canvas.width
        if (i === 0) ctx.moveTo(x, y)
        else ctx.lineTo(x, y)
      }
      ctx.stroke()
      break

    case 'circular':
      const centerX = canvas.width / 2
      const centerY = canvas.height / 2
      const radius = Math.min(canvas.width, canvas.height) / 3
      for (let i = 0; i < bufferLength; i++) {
        const angle = (i / bufferLength) * Math.PI * 2
        const amplitude = (dataArray[i] / 255) * radius
        const x = centerX + Math.cos(angle) * (radius + amplitude)
        const y = centerY + Math.sin(angle) * (radius + amplitude)
        ctx.beginPath()
        ctx.arc(x, y, amplitude / 4, 0, Math.PI * 2)
        ctx.fillStyle = `hsl(${(i * 360) / bufferLength}, 70%, 60%)`
        ctx.fill()
      }
      break

    case 'particles':
      for (let i = 0; i < 100; i++) {
        const x = Math.random() * canvas.width
        const y = Math.random() * canvas.height
        const size = (dataArray[Math.floor(Math.random() * dataArray.length)] / 255) * 8
        ctx.fillStyle = `rgba(102, 126, 234, ${size / 8})`
        ctx.fillRect(x, y, size, size)
      }
      break

    case 'fire':
      for (let i = 0; i < bufferLength; i++) {
        const intensity = dataArray[i] / 255
        const barHeight = intensity * canvas.height
        const gradient = ctx.createLinearGradient(0, canvas.height - barHeight, 0, canvas.height)
        gradient.addColorStop(0, '#ff4500')
        gradient.addColorStop(0.5, '#ff6600')
        gradient.addColorStop(1, '#ffaa00')
        ctx.fillStyle = gradient
        ctx.fillRect(i * barWidth, canvas.height - barHeight, barWidth - 1, barHeight)
      }
      break
  }
}

// ========== Запись ==========
async function startRecording() {
  try {
    mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true })

    // Пробуем разные MIME типы
    let mimeType = ''
    const types = ['audio/webm', 'audio/mp4', 'audio/ogg']
    for (const type of types) {
      if (MediaRecorder.isTypeSupported(type)) {
        mimeType = type
        break
      }
    }

    mediaRecorder = new MediaRecorder(mediaStream, { mimeType })
    audioChunks = []

    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.push(event.data)
      }
    }

    mediaRecorder.onstop = async () => {
      const blob = new Blob(audioChunks, { type: mimeType || 'audio/webm' })
      audioBlob.value = blob

      // Создаём URL для аудио
      if (audioUrl.value) {
        URL.revokeObjectURL(audioUrl.value)
      }
      audioUrl.value = URL.createObjectURL(blob)

      // Генерируем волну
      await generateWaveformFromBlob(blob)

      // Создаём аудио элемент
      if (audioElement) {
        audioElement.pause()
        audioElement.src = audioUrl.value
      }

      // Останавливаем визуализацию
      stopVisualization()

      if (durationInterval) {
        clearInterval(durationInterval)
        durationInterval = null
      }

      // Закрываем стрим
      if (mediaStream) {
        mediaStream.getTracks().forEach((track) => track.stop())
        mediaStream = null
      }
    }

    mediaRecorder.start(100)
    isRecording.value = true

    // Настройка визуализации записи
    recordingContext = new AudioContext()
    analyserNode = recordingContext.createAnalyser()
    analyserNode.fftSize = 256
    sourceNode = recordingContext.createMediaStreamSource(mediaStream)
    sourceNode.connect(analyserNode)

    const bufferLength = analyserNode.frequencyBinCount
    const dataArray = new Uint8Array(bufferLength)

    function animate() {
      if (!isRecording.value) return
      animationId = requestAnimationFrame(animate)
      analyserNode?.getByteFrequencyData(dataArray)
      drawRecordingVisualization(dataArray)
    }

    animate()

    recordingDuration.value = 0
    durationInterval = window.setInterval(() => {
      recordingDuration.value++
    }, 1000)
  } catch (error) {
    console.error('Ошибка:', error)
    alert('Не удалось получить доступ к микрофону')
  }
}

function stopRecording() {
  if (mediaRecorder && mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
    isRecording.value = false
  }
}

function stopVisualization() {
  if (animationId) {
    cancelAnimationFrame(animationId)
    animationId = null
  }
  if (recordingContext) {
    recordingContext.close()
    recordingContext = null
  }
}

function toggleRecording() {
  if (isRecording.value) {
    stopRecording()
  } else if (audioUrl.value) {
    resetRecording()
    startRecording()
  } else {
    startRecording()
  }
}

// ========== Воспроизведение (через HTMLAudioElement) ==========
function togglePlayback() {
  if (!audioUrl.value) return

  if (!audioElement) {
    audioElement = new Audio(audioUrl.value)
    setupAudioListeners()
  }

  if (isPlaying.value) {
    audioElement.pause()
    if (progressUpdateInterval) {
      clearInterval(progressUpdateInterval)
      progressUpdateInterval = null
    }
    isPlaying.value = false
  } else {
    audioElement.play()
    isPlaying.value = true

    // Обновляем прогресс
    if (progressUpdateInterval) clearInterval(progressUpdateInterval)
    progressUpdateInterval = setInterval(() => {
      if (audioElement && isPlaying.value) {
        currentTime.value = audioElement.currentTime
        playbackProgress.value = (currentTime.value / duration.value) * 100

        if (currentTime.value >= duration.value) {
          stopPlayback()
        }
      }
    }, 50)
  }
}

function setupAudioListeners() {
  if (!audioElement) return

  audioElement.addEventListener('ended', () => {
    stopPlayback()
  })

  audioElement.addEventListener('error', (e) => {
    console.error('Ошибка аудио:', e)
    alert('Ошибка воспроизведения')
    stopPlayback()
  })

  audioElement.addEventListener('canplay', () => {
    duration.value = audioElement!.duration
  })
}

function stopPlayback() {
  if (audioElement) {
    audioElement.pause()
    audioElement.currentTime = 0
  }
  if (progressUpdateInterval) {
    clearInterval(progressUpdateInterval)
    progressUpdateInterval = null
  }
  isPlaying.value = false
  currentTime.value = 0
  playbackProgress.value = 0
}

function seekAudio(event: MouseEvent) {
  if (!waveformCanvasRef.value || !audioUrl.value) return

  const rect = waveformCanvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const percent = Math.max(0, Math.min(1, x / rect.width))
  const newTime = percent * duration.value

  if (audioElement) {
    const wasPlaying = isPlaying.value
    audioElement.currentTime = newTime
    currentTime.value = newTime
    playbackProgress.value = percent * 100

    if (!wasPlaying) {
      // Если не играло, просто обновляем позицию
      pausedAt = newTime
    }
  }
}

// Для seekAudio
let pausedAt = 0

function resetRecording() {
  stopPlayback()
  if (audioElement) {
    audioElement.pause()
    audioElement.src = ''
  }
  if (audioUrl.value) {
    URL.revokeObjectURL(audioUrl.value)
  }
  audioUrl.value = ''
  audioBlob.value = null
  recordingDuration.value = 0
  duration.value = 0
  currentTime.value = 0
  playbackProgress.value = 0
  waveformData = []

  if (waveformCanvasRef.value) {
    const ctx = waveformCanvasRef.value.getContext('2d')
    if (ctx) {
      ctx.fillStyle = '#0a0e27'
      ctx.fillRect(0, 0, waveformCanvasRef.value.width, waveformCanvasRef.value.height)
    }
  }
}

function downloadAudio() {
  if (audioBlob.value) {
    const url = URL.createObjectURL(audioBlob.value)
    const a = document.createElement('a')
    a.href = url
    a.download = `recording_${new Date().toISOString()}.webm`
    a.click()
    URL.revokeObjectURL(url)
  }
}

// ========== Очистка ==========
onUnmounted(() => {
  if (durationInterval) clearInterval(durationInterval)
  if (progressUpdateInterval) clearInterval(progressUpdateInterval)
  if (animationId) cancelAnimationFrame(animationId)
  if (audioElement) {
    audioElement.pause()
    audioElement.src = ''
  }
  if (audioUrl.value) URL.revokeObjectURL(audioUrl.value)
  if (recordingContext) recordingContext.close()
  if (mediaStream) {
    mediaStream.getTracks().forEach((track) => track.stop())
  }
})

// Ресайз
import { onMounted } from 'vue'
onMounted(() => {
  window.addEventListener('resize', () => {
    if (waveformData.length) drawWaveform()
  })
})
</script>

<style scoped>
/* Стили такие же как в предыдущей версии, они не изменились */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-family:
    'Inter',
    system-ui,
    -apple-system,
    sans-serif;
  padding: 20px;
}

.header {
  text-align: center;
  color: white;
  margin-bottom: 30px;
}

.title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 10px;
}

.subtitle {
  font-size: 18px;
  opacity: 0.9;
}

.style-selector {
  max-width: 800px;
  margin: 0 auto 30px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 20px;
}

.selector-title {
  color: white;
  text-align: center;
  margin-bottom: 15px;
  font-size: 18px;
}

.style-buttons {
  display: flex;
  gap: 10px;
  justify-content: center;
  flex-wrap: wrap;
}

.style-btn {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 30px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.style-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.style-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: white;
}

.recorder-area {
  max-width: 800px;
  margin: 0 auto 40px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 30px;
  padding: 30px;
}

.visualizer-container {
  position: relative;
  background: #0a0e27;
  border-radius: 20px;
  overflow: hidden;
  margin-bottom: 30px;
}

.visualizer-canvas {
  width: 100%;
  height: 300px;
  display: block;
}

.visualizer-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(10, 14, 39, 0.8);
  color: white;
  font-size: 18px;
}

.record-btn {
  position: relative;
  display: block;
  margin: 0 auto 20px;
  padding: 18px 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 60px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.record-btn.recording {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  animation: pulse 1.5s infinite;
}

.btn-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.recording-pulse {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 60px;
  animation: ripple 1s infinite;
}

.timer {
  text-align: center;
}

.timer-circle {
  display: inline-flex;
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  align-items: center;
  justify-content: center;
}

.timer-text {
  font-size: 28px;
  font-weight: 700;
  color: white;
  font-family: monospace;
}

.playback-area {
  max-width: 800px;
  margin: 0 auto;
}

.playback-card {
  background: white;
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.playback-card h3 {
  color: #333;
  margin-bottom: 20px;
  font-size: 24px;
  text-align: center;
}

.waveform-container {
  position: relative;
  background: #0a0e27;
  border-radius: 15px;
  overflow: hidden;
  margin-bottom: 20px;
  cursor: pointer;
}

.waveform-canvas {
  width: 100%;
  height: 150px;
  display: block;
  cursor: pointer;
}

.progress-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
}

.progress-bar {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: rgba(102, 126, 234, 0.3);
  transition: width 0.1s linear;
  pointer-events: none;
}

.playhead {
  position: absolute;
  top: 0;
  height: 100%;
  width: 2px;
  pointer-events: none;
  transition: left 0.1s linear;
}

.playhead-dot {
  position: absolute;
  top: 50%;
  left: -6px;
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 50%;
  transform: translateY(-50%);
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.8);
}

.custom-controls {
  display: flex;
  gap: 15px;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
}

.play-btn {
  padding: 12px 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 50px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: transform 0.2s;
}

.play-btn:hover {
  transform: scale(1.05);
}

.time-info {
  display: flex;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #667eea;
  font-family: monospace;
}

.download-btn,
.reset-btn {
  padding: 12px 24px;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.download-btn {
  background: #42b883;
  color: white;
}

.reset-btn {
  background: #f0f0f0;
  color: #666;
}

.download-btn:hover,
.reset-btn:hover {
  transform: translateY(-2px);
}

.tips {
  max-width: 600px;
  margin: 40px auto 0;
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}

.tip-item {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(5px);
  padding: 10px 20px;
  border-radius: 30px;
  color: white;
  font-size: 14px;
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

@keyframes ripple {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

@media (max-width: 768px) {
  .title {
    font-size: 32px;
  }
  .visualizer-canvas {
    height: 200px;
  }
  .waveform-canvas {
    height: 100px;
  }
  .custom-controls {
    flex-direction: column;
  }
  .download-btn,
  .reset-btn {
    width: 100%;
  }
}
</style>
