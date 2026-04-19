import multiprocessing
import logging
import sys

from snaker.app.config import kafka_address
from snaker.app.infrastructure.kafka.adapter.audio_adapter import AudioWorker
from snaker.app.infrastructure.kafka.adapter.llm_adapter import LLMWorker
from snaker.app.infrastructure.kafka.adapter.output_adapter import OutputWorker
from snaker.app.infrastructure.kafka.adapter.pii_router_worker import PIIRoutingWorker
from snaker.app.infrastructure.kafka.adapter.stt_adapter import STTWorker
from snaker.app.infrastructure.kafka.adapter.times_adapter import TimesWorker

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(name)s: %(message)s',
    handlers=[logging.StreamHandler(sys.stdout)]
)

def run_worker(worker_class, bootstrap_servers):
    """Функция-обертка для запуска воркера в процессе"""
    try:
        worker = worker_class(bootstrap_servers=bootstrap_servers)
        worker.start()
    except Exception as e:
        logging.error(f"Fatal error in {worker_class.__name__}: {e}")

if __name__ == "__main__":
    BOOTSTRAP_SERVERS = f"{kafka_address}" # Или адрес из конфига

    # Список воркеров для запуска
    workers = [
        STTWorker,
        PIIRoutingWorker,
        LLMWorker,
        TimesWorker,
        AudioWorker,
        OutputWorker
    ]

    processes = []

    logging.info("Starting 3dom-pipeline workers...")

    for worker_cls in workers:
        p = multiprocessing.Process(
            target=run_worker,
            args=(worker_cls, BOOTSTRAP_SERVERS),
            name=worker_cls.__name__
        )
        p.start()
        processes.append(p)
        logging.info(f"Started process for {worker_cls.__name__} (PID: {p.pid})")

    try:
        for p in processes:
            p.join()
    except KeyboardInterrupt:
        logging.info("Shutdown signal received. Terminating workers...")
        for p in processes:
            p.terminate()
            p.join()
