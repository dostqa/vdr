from pydantic import BaseModel

class ObjectPdn(BaseModel):
    text: str
    type: str
    start_time: float
    end_time: float

class AudioMessage(BaseModel):
    request_id: int
    file_path: str
    original_text: str
    anon_text: str
    objects_pdns: list[ObjectPdn]
