
from rapidfuzz import fuzz
import re

from snaker.app.domain.message.stt import WordData


def normalize(text: str) -> list[str]:
    text = text.lower()
    text = re.sub(r"[^\w\s]", "", text)
    return text.split()

def find_span(obj_text: str, whisper_words: list[WordData]):
    obj_tokens = normalize(obj_text)
    if not obj_tokens:
        return None

    whisper_tokens = [normalize(w.text)[0] for w in whisper_words if normalize(w.text)]

    best_score = 0
    best_span = None

    window_size = len(obj_tokens)

    for i in range(len(whisper_tokens) - window_size + 1):
        window = whisper_tokens[i:i + window_size]

        score = fuzz.ratio(
            " ".join(obj_tokens),
            " ".join(window)
        )

        if score > best_score:
            best_score = score
            best_span = (i, i + window_size - 1)

    if best_span:
        start, end = best_span
        return whisper_words[start].start_time, whisper_words[end].end_time

    return None
