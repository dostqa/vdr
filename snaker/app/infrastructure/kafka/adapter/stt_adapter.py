from snaker.app.domain.message.input import InputMessage
from snaker.app.domain.message.stt import STTMessage
from snaker.app.infrastructure.kafka.kafka_transport import BaseKafkaWorker
from snaker.app.infrastructure.minio.minio_repository import MinioRepo
from snaker.app.service.whisper_service import WhisperService


class STTWorker(BaseKafkaWorker[InputMessage, STTMessage]):
    def __init__(self, bootstrap_servers: str):
        super().__init__(
            bootstrap_servers=bootstrap_servers,
            group_id="stt_service_group",
            in_topic="input_topic",
            out_topic="stt_topic",
            in_model=InputMessage
        )

        self.minio = MinioRepo()
        self.engine = WhisperService()

    def handle(self, data: InputMessage) -> STTMessage:
        print(f"[*] STT (ID: {data.request_id})")

        self.minio.get(data.file_path, f"./tmp/{data.file_path}")

        full_text, words = self.engine.transcription(f"./tmp/{data.file_path}")

        return STTMessage(
            request_id=data.request_id,
            file_path=data.file_path,
            full_text=full_text,
            words=words
        )