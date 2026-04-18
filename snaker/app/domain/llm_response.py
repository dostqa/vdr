from pydantic import BaseModel
from typing import List

from snaker.app.domain.message.pers import TypesPers


class Object(BaseModel):
    raw_text: str
    clean_text: str
    type: TypesPers

class LLMResponse(BaseModel):
    original_text: str
    anon_text: str
    objects: list[Object]
