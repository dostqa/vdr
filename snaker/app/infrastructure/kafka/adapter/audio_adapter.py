from snaker.app.domain.message.audio import AudioMessage
from snaker.app.domain.message.times import TimesMessage
from snaker.app.infrastructure.kafka.kafka_transport import BaseKafkaWorker
from snaker.app.service.audio_service import redact_audio_pydub


class AudioWorker(BaseKafkaWorker[TimesMessage, AudioMessage]):
    def __init__(self, bootstrap_servers: str):
        super().__init__(
            bootstrap_servers=bootstrap_servers,
            group_id="stt_service_group",
            in_topic="times_topic",
            out_topic="audio_topic",
            in_model=AudioMessage
        )

    def handle(self, data: TimesMessage) -> AudioMessage:
        print(f"[*] Audio (ID: {data.request_id})")


        redact_audio_pydub(input_path=f"./tmp/{data.file_path}", output_path=f"./tmp/anon_{data.file_path}", pdns=data.objects_pdns)

        return AudioMessage(
            request_id=data.request_id,
            file_path=data.file_path,
            original_text=data.original_text,
            anon_text=data.anon_text,
            objects_pdns=data.objects_pdns,
        )