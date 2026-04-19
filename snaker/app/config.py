import os

minio_user = os.environ.get('MINIO_ROOT_USER', 'admin')
minio_password = os.environ.get('MINIO_ROOT_PASSWORD', "adminpass")
minio_address="localhost:9000"
hf_token = "478233b0b7a84813b6ab9b5a86bcc761.tJvUMRJJMCtVnL6VeMKIzF9L"
kafka_address = "localhost:9092"

