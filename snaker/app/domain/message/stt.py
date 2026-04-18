import time

from pydantic import BaseModel

class WordData(BaseModel):
    text: str
    start_time: float
    end_time: float

class STTMessage(BaseModel):
    request_id: int
    file_path: str
    full_text: str
    words: list[WordData]