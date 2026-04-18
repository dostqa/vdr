from pydantic import BaseModel

from snaker.app.domain.message.audio import ObjectPdn


class TimesMessage(BaseModel):
    request_id: int
    file_path: str
    original_text: str
    anon_text: str
    objects_pdns: list[ObjectPdn]
