from snaker.app.domain.message.llm import LLMMessage, ObjectLLM
from snaker.app.domain.message.stt import STTMessage
from snaker.app.infrastructure.kafka.kafka_transport import BaseKafkaWorker
from snaker.app.service.llm_service import LLMService


class LLMWorker(BaseKafkaWorker[STTMessage, LLMMessage]):
    def __init__(self, bootstrap_servers: str):
        super().__init__(
            bootstrap_servers=bootstrap_servers,
            group_id="stt_service_group",
            in_topic="check_topic",
            out_topic="llm_topic",
            in_model=STTMessage
        )

        self.llm = LLMService()

    def handle(self, data: STTMessage) -> LLMMessage:
        print(f"[*] LLM (ID: {data.request_id})")

        llm_response = self.llm.gen(data.full_text)
        objects_llm = []

        for obj in llm_response.objects:
            obj_llm = ObjectLLM(
                raw_text=obj.raw_text,
                clean_text=obj.clean_text,
                type=obj.type
            )
            objects_llm.append(obj_llm)

        return LLMMessage(
            request_id=data.request_id,
            file_path=data.file_path,
            original_text=llm_response.original_text,
            anon_text=llm_response.anon_text,
            words=data.words,
            objects_llm=objects_llm
        )