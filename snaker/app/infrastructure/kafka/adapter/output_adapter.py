import os

from snaker.app.domain.message.audio import AudioMessage
from snaker.app.domain.message.output import OutputMessage
from snaker.app.infrastructure.kafka.kafka_transport import BaseKafkaWorker
from snaker.app.infrastructure.minio.minio_repository import MinioRepo


class OutputWorker(BaseKafkaWorker[AudioMessage, OutputMessage]):
    def __init__(self, bootstrap_servers: str):
        super().__init__(
            bootstrap_servers=bootstrap_servers,
            group_id="stt_service_group",
            in_topic="audio_topic",
            out_topic="output_topic",
            in_model=AudioMessage
        )

        self.minio = MinioRepo()

    def handle(self, data: AudioMessage) -> OutputMessage:
        print(f"[*] Output (ID: {data.request_id})")

        self.minio.save(f"anon_{data.file_path}", f"./tmp/anon_{data.file_path}")

        os.remove(f"./tmp/anon_{data.file_path}")
        os.remove(f"./tmp/{data.file_path}")

        return OutputMessage(
            request_id=data.request_id,
            old_file_path=data.file_path,
            new_file_path=f"anon_{data.file_path}",
            original_text=data.original_text,
            anon_text=data.anon_text,
            objects_pdns=data.objects_pdns,
        )