import re
from natasha import Segmenter, NewsEmbedding, NewsNERTagger, Doc

class PIIIdentifier:
    def __init__(self):
        self.segmenter = Segmenter()
        self.ner_tagger = NewsNERTagger(NewsEmbedding())
        # Компактный список паттернов: Телефон, Email, ИНН/СНИЛС/Паспорт
        self.regex = re.compile(
            r'(\+7|8)[\s-]?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}|' # Телефон
            r'[\w\.-]+@[\w\.-]+\.\w{2,4}|'                                # Email
            r'\b\d{10}\b|\b\d{12}\b|'                                     # ИНН
            r'\b\d{3}-\d{3}-\d{3}\s\d{2}\b|'                              # СНИЛС
            r'\b\d{4}\s\d{6}\b'                                           # Паспорт
        )

    def has_pii(self, text: str) -> bool:
        # 1. Быстрая проверка регулярками
        if self.regex.search(text):
            return True

        # 2. Проверка имен и адресов через Natasha
        doc = Doc(text)
        doc.segment(self.segmenter)
        doc.tag_ner(self.ner_tagger)

        return any(span.type in ['PER', 'LOC'] for span in doc.spans)
