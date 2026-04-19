# Architcture
1) Go получает файл от Vue.js сохраняет его в Minio
2) Отправляет в InputTopic сообщение об этом
3) Python Worker начиобработку
4) STT (fast_whisper) преорабзует аудио в текст и передает в PIIChecker
5) С помощью reg и NER(Natasha) проверяем наличие ПДн
6) Если есть, то отправляет в LLM, если нет то отправляем в Output Topic
7) LLM выделяет ПДн
8) RappidFuzz совмещяет ПДн с временными таймингами
9) Глушится изначальное видео в нужных местах
10) Сохраняется в Minio и отдается на Vue js для скачивания

# Запуск
```
docker-compose up --build
```
P.S. при необходимости поменять конфигурацию в snake/config.py или gopher/configs/config_docker.yaml
