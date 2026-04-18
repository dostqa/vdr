from pydub import AudioSegment
from typing import List

from snaker.app.domain.message.audio import ObjectPdn

def redact_audio_pydub(input_path: str, output_path: str, pdns: List[ObjectPdn]):
    audio = AudioSegment.from_file(input_path)

    pdns.sort(key=lambda x: x.start_time)

    for pdn in pdns:
        start_ms = int(pdn.start_time * 1000)
        end_ms = int(pdn.end_time * 1000)
        duration = end_ms - start_ms

        quiet = AudioSegment.silent(duration=duration)

        audio = audio[:start_ms] + quiet + audio[end_ms:]

    audio.export(output_path, format="webm")
    print(f"Файл сохранен: {output_path}")