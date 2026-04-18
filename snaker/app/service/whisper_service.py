from typing import Tuple

from loguru import logger
from faster_whisper import WhisperModel

from snaker.app.domain.message.stt import WordData
from snaker.app.service.singelton import Singleton


class WhisperService(metaclass=Singleton):
    def __init__(self, model_size="small", device="cpu", compute_type="int8"):
        logger.info("Creating Whisper Service...")
        self.model = WhisperModel(model_size, device=device, compute_type=compute_type)

    def transcription(self, filename: str, language="ru", batch_size=16) -> Tuple[str, list[WordData]]:
        logger.info(f"Start transcription for file: {filename}")
        segments, _ = self.model.transcribe(filename, word_timestamps=True)
        words = []
        full_text = ""
        for segment in segments:
            full_text += segment.text
            print(segment)
            for word in segment.words:
               word_model = WordData(text=word.word, start_time=word.start, end_time=word.end)
               words.append(word_model)

        return (full_text, words)