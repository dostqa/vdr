import logging
from typing import List

from snaker.app.infrastructure.kafka.adapter.checker_adapter import PIIRouter
from snaker.app.service.piichecker import PIIIdentifier


# Класс-адаптер для встраивания в multiprocessing pipeline
class PIIRoutingWorker:
    def __init__(self, bootstrap_servers: str):
        self.bootstrap_servers = bootstrap_servers
        self.group_id = "pii_router_group"

        # Названия топиков согласно вашей архитектуре
        self.in_topic = "stt_topic"
        self.check_topic = "check_topic"
        self.output_topic = "output_topic"

        # Инициализируем детектор один раз при создании воркера
        self.checker = PIIIdentifier()

        # Инициализируем сам роутер
        self.router = PIIRouter(
            bootstrap_servers=self.bootstrap_servers,
            group_id=self.group_id,
            pii_checker=self.checker
        )

    def start(self):
        """Этот метод вызовет run_worker"""
        logging.info(f"PIIRoutingWorker started logic: {self.in_topic} -> [{self.check_topic} | {self.output_topic}]")
        self.router.run(
            in_topic=self.in_topic,
            check_topic=self.check_topic,
            output_topic=self.output_topic
        )