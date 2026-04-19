# Используем базовый образ Python
FROM python:3.11-slim

# Ставим рабочую директорию
WORKDIR /app

# Копируем список зависимостей
COPY requirements.txt .

# Устанавливаем библиотеки
RUN pip install pydantic confluent_kafka minio loguru pydub ollama rapidfuzz faster_whisper

# Копируем все остальные файлы проекта
COPY . .

# Запускаем скрипт напрямую
CMD ["python", "main.py"]