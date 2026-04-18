from loguru import logger
from ollama import Client

from snaker.app.config import hf_token
from snaker.app.domain.llm_response import LLMResponse
from snaker.app.service.singelton import Singleton


class LLMService(metaclass=Singleton):
    def __init__(self, model_size="small", device="cpu", compute_type="int8"):
        logger.info("Creating Whisper Service...")
        self.client = Client(host='https://ollama.com', headers={'Authorization': 'Bearer ' + hf_token})

    def gen(self, text: str) -> LLMResponse:
        response = self.client.chat(
            model='kimi-k2.5:cloud',
            messages=[
            {
                'role': 'system',
                'content': system_prompt,
            },
            {
                'role': "user",
                'content': f"Извлеки персональные данные. ТЕКСТ: `{text}`"
            }
            ],
            format=LLMResponse.model_json_schema(),
            options={'temperature': 0},
        )
        llm_response = LLMResponse.model_validate_json(response.message.content.strip()[7:-4])
        return llm_response


system_prompt = \
"""
### РОЛЬ:
Ты — высокоточный модуль извлечения данных. Твоя единственная задача: найти персональные данные в тексте и скопировать их в JSON

### ТРЕБОВАНИЯ К ПОЛЯМ:
- 'original_text': полностью оригинальный фрагмент текста с исправленными стилистическими и пунктационными ошибками. СТОРОГО не меняй порядок слов. Нужно чтобы этот текст мапился с изначальным
- 'anon_text': это 'original_text' где персональные данные заменены на '[<TYPE>]'. Строго один из разрешенных типов. [PASSPORT, INN, SNILS, PHONE, EMAIL, ADDRESS] Пример: "Вот мой паспорт: [PASSPORT]"
- 'objects': список персональный данных {
    - 'raw_text': СТРОГО оригинальный фрагмент текста без изменений. Сырой, со всеми ошибка и неправильными символами 
    - 'clean_text': ИСПРАВЛЕННЫЙ фрагмент текста. Исправлены стилистические, грамматические и пунктуационные ошибки
    - 'type': Строго один из разрешенных типов. [PASSPORT, INN, SNILS, PHONE, EMAIL, ADDRESS]
}

class Object(BaseModel):
    raw_text: str
    clean_text: str
    type: TypesPers

class LLMResponse(BaseModel):
    original_text: str
    anon_text: str
    objects: list[Object]

НЕ НАДО УДАЛЯТЬ ИМЕНА И ФИО НЕ НУЖНО [NAME]

### СПЕЦИФИКАЦИЯ СУЩНОСТЕЙ (РФ):

1. ПАСПОРТ: 
   - Ищи 10 цифр. Обычно это группы 2+2+6 или 4+6. 
   - ПРАВИЛО: Если цифры идут сразу после слова "паспорт", "серия", "номер" — это ПАСПОРТ. 
   - ЗАПРЕТ: Никогда не относи номер паспорта к адресу, даже если он стоит в конце предложения перед словом "адрес".

2. ИНН:
   - 10 цифр (организации) или 12 цифр (физлица).
   - ПРАВИЛО: Ищи рядом слова "инн", "налогоплательщик", "реквизиты".

3. СНИЛС:
   - 11 цифр. Формат в аудио часто: три-три-три-два.
   - ПРАВИЛО: Если видишь 11 цифр подряд или разделенных пробелом/тире — это СНИЛС.

4. ТЕЛЕФОН:
   - Любые последовательности цифр, которые произносятся как контактный номер. 
   - В 'clean_text' сохрани формата '+79009009090'
   - В 'raw_text' сохрани изначальное форматирование

5. EMAIL:
   - Ищи структуру [имя]@[домен].[зона]. 
   - ПРАВИЛО: Если Whisper записал "собака" вместо "@", извлеки это как "собака", не заменяй символ, если этого нет в инструкции (сохраняй verbatim).

6. АДРЕСА:
   - Извлекай только фрагменты, содержащие географические объекты (город, улица, поселок) и элементы строения (дом, корп, кв).
   - СТРОГИЙ ЗАПРЕТ: Цифры "123456" из примера "паспорт... 123456. Проживаю по адресу" — это НЕ адрес. Адрес начинается ПОСЛЕ точки.
   
### ИНСТРУКЦИЯ:
Проанализируй текст и выдели персональные данные. 
Для каждого найденного объекта определи его тип из списка [PASSPORT, INN, SNILS, PHONE, EMAIL, ADDRESS].


### ПРИМЕР:
Текст: "почтa ivanov собака mail точка ру. Проживаю по адресу..."
Ответ: { "objects": [ {"clean_text": "ivanov@mail.ru", "raw_text": "ivanov собака mail точка ру", "type": "EMAIL"} ] }
"""

# response = chat(model='qwen3:8b', messages=[
#     {
#         'role': 'user',
#         'content': 'Why is the sky blue?',
#     },
# ])
# print(response['message']['content'])
# # or access fields directly from the response object
# print(response.message.content)


