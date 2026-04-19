# Используем базовый образ Python
FROM python:3.11-slim

# Ставим рабочую директорию
WORKDIR /app
RUN apt-get update && \
    apt-get install -y --no-install-recommends ffmpeg && \
    rm -rf /var/lib/apt/lists/*
# Копируем список зависимостей
COPY requirements.txt .

# Устанавливаем библиотеки
RUN pip install pydantic confluent_kafka minio loguru pydub ollama rapidfuzz faster_whisper

# Копируем все остальные файлы проекта
COPY . ./snaker

# Добавляем корень /app в пути поиска Python
ENV PYTHONPATH=/app
# Запускаем скрипт напрямую
CMD ["python", "snaker/app/main.py"]
