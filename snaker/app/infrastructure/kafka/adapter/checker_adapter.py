import logging
from pydantic import BaseModel
from typing import List, Optional
from confluent_kafka import Consumer, Producer

from snaker.app.domain.message.output import OutputMessage
from snaker.app.domain.message.stt import STTMessage
from snaker.app.service.piichecker import PIIIdentifier


class PIIRouter:
    def __init__(
            self,
            bootstrap_servers: str,
            group_id: str,
            pii_checker: PIIIdentifier
    ):
        self.pii_checker = pii_checker
        self.consumer = Consumer({
            'bootstrap.servers': bootstrap_servers,
            'group.id': group_id,
            'auto.offset.reset': 'earliest',
            'enable.auto.commit': False
        })
        self.producer = Producer({'bootstrap.servers': bootstrap_servers})

    def send(self, topic: str, model: BaseModel):
        """Сериализация и отправка в Kafka"""
        payload = model.model_dump_json().encode('utf-8')
        # Ключ по request_id для консистентности в партициях
        self.producer.produce(
            topic,
            key=str(model.request_id),
            value=payload
        )
        self.producer.flush()

    def run(self, in_topic: str, check_topic: str, output_topic: str):
        logging.info(f"PII Router active. Listening: {in_topic}")
        self.consumer.subscribe([in_topic])

        try:
            while True:
                msg = self.consumer.poll(1.0)
                if msg is None: continue
                if msg.error():
                    logging.error(f"Kafka error: {msg.error()}")
                    continue

                try:
                    # 1. Получаем входящее сообщение (STTMessage)
                    raw_value = msg.value().decode('utf-8')
                    data = STTMessage.model_validate_json(raw_value)

                    # 2. Проверяем текст на наличие ПДн
                    has_pii = self.pii_checker.has_pii(data.full_text)

                    if has_pii:
                        # ЕСТЬ ПДН -> Отправляем в check_topic (на доработку LLM)
                        # Возвращаем ту же структуру STTMessage
                        self.send(check_topic, data)
                        logging.info(f"[{data.request_id}] PII found -> Sent to {check_topic}")
                    else:
                        # НЕТ ПДН -> Сразу формируем финальный OutputMessage
                        out_msg = OutputMessage(
                            request_id=data.request_id,
                            old_file_path=data.file_path,
                            new_file_path=data.file_path, # Файл не менялся
                            original_text=data.full_text,
                            anon_text=data.full_text,     # Текст не менялся
                            objects_pdns=[]               # Список пуст
                        )
                        self.send(output_topic, out_msg)
                        logging.info(f"[{data.request_id}] Clean -> Sent to {output_topic}")

                    # Фиксируем смещение (offset) в Kafka
                    self.consumer.commit(asynchronous=False)

                except Exception as e:
                    logging.error(f"Error processing {msg.key()}: {e}")

        except KeyboardInterrupt:
            logging.info("Service stopping...")
        finally:
            self.consumer.close()
