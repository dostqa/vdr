from pydantic import BaseModel

from snaker.app.domain.message.pers import TypesPers
from snaker.app.domain.message.stt import WordData


class ObjectLLM(BaseModel):
    raw_text: str
    clean_text: str
    type: TypesPers

class LLMMessage(BaseModel):
    request_id: int
    file_path: str
    original_text: str
    anon_text: str
    words: list[WordData]
    objects_llm: list[ObjectLLM]
