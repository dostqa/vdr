from confluent_kafka import Producer
from minio import Minio

from snaker.app.domain.message.input import InputMessage
client = Minio(
    "localhost:9000",
    access_key="admin",
    secret_key="adminpass",
    secure=False
)

bucket_name = "bucket"
file_path = "audio.mp3" # give correct file path
object_name = "audio.mp3"  # File name in the bucket

client.fput_object(bucket_name, object_name, file_path)

input()
# 1. Настройка
conf = {'bootstrap.servers': 'localhost:9092'}
producer = Producer(conf)

# 2. Подготовка данных
message = InputMessage(request_id=123, file_path="audio.mp3")
payload = message.model_dump_json().encode('utf-8')

# 3. Отправка
producer.produce(
    topic='input_topic',
    value=payload
)

# 4. Ждем, пока сообщение реально уйдет в сеть
producer.flush()