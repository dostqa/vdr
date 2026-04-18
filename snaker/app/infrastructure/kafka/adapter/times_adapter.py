from snaker.app.domain.message.llm import LLMMessage
from snaker.app.domain.message.audio import AudioMessage, ObjectPdn
from snaker.app.domain.message.times import TimesMessage
from snaker.app.infrastructure.kafka.kafka_transport import BaseKafkaWorker
from snaker.app.service.times_service import find_span


class TimesWorker(BaseKafkaWorker[LLMMessage, TimesMessage]):
    def __init__(self, bootstrap_servers: str):
        super().__init__(
            bootstrap_servers=bootstrap_servers,
            group_id="stt_service_group",
            in_topic="llm_topic",
            out_topic="times_topic",
            in_model=LLMMessage
        )

    def handle(self, data: LLMMessage) -> TimesMessage:
        print(f"[*] Times (ID: {data.request_id})")

        objects_pdns = []

        for obj in data.objects_llm:
            span = find_span(obj.raw_text, data.words)

            if span:
                start_time, end_time = span
            else:
                start_time, end_time = 0.0, 0.0

            objects_pdns.append(
                ObjectPdn(
                    text=obj.clean_text,
                    type=obj.type,
                    start_time=start_time,
                    end_time=end_time
                )
            )

        return TimesMessage(
            request_id=data.request_id,
            file_path=data.file_path,
            original_text=data.original_text,
            anon_text=data.anon_text,
            objects_pdns=objects_pdns
        )