import logging
from typing import Type, TypeVar, Generic
from pydantic import BaseModel
from confluent_kafka import Consumer, Producer

T_IN = TypeVar("T_IN", bound=BaseModel)
T_OUT = TypeVar("T_OUT", bound=BaseModel)

class BaseKafkaWorker(Generic[T_IN, T_OUT]):
    def __init__(
            self,
            bootstrap_servers: str,
            group_id: str,
            in_topic: str,
            out_topic: str,
            in_model: Type[T_IN]
    ):
        self.in_topic = in_topic
        self.out_topic = out_topic
        self.in_model = in_model

        self.consumer = Consumer({
            'bootstrap.servers': bootstrap_servers,
            'group.id': group_id,
            'auto.offset.reset': 'earliest',
            'enable.auto.commit': False,
            'max.poll.interval.ms': 600000  # 10 минут на обработку файла
        })

        self.producer = Producer({'bootstrap.servers': bootstrap_servers})
        self.consumer.subscribe([self.in_topic])

    def start(self):
        logging.info(f"Worker started. Listening to {self.in_topic}...")
        try:
            while True:
                msg = self.consumer.poll(1.0)
                if msg is None: continue
                if msg.error():
                    logging.error(f"Kafka error: {msg.error()}")
                    continue

                try:
                    raw_data = msg.value().decode('utf-8')
                    input_data = self.in_model.model_validate_json(raw_data)

                    output_data = self.handle(input_data)

                    if output_data:
                        self.send(output_data)
                        self.consumer.commit(asynchronous=False)

                except Exception as e:
                    logging.exception(f"Error processing message: {e}")
        finally:
            self.consumer.close()

    def send(self, model: T_OUT):
        payload = model.model_dump_json().encode('utf-8')
        self.producer.produce(self.out_topic, key=str(model.request_id), value=payload)
        self.producer.flush()

    def handle(self, data: T_IN) -> T_OUT:
        raise NotImplementedError