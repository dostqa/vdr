from pydantic import BaseModel

from snaker.app.domain.message.audio import ObjectPdn


class OutputMessage(BaseModel):
    request_id: int
    old_file_path: str
    new_file_path: str
    original_text: str
    anon_text: str
    objects_pdns: list[ObjectPdn]