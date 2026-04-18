from pydantic import BaseModel

class InputMessage(BaseModel):
    request_id: int
    file_path: str
