from loguru import logger
from minio import Minio

from snaker.app.config import minio_user, minio_password
from snaker.app.infrastructure.minio.singelton import Singleton


class MinioRepo(metaclass=Singleton):
    def __init__(self, model_size="small", device="cpu", compute_type="int8"):
        logger.info("Creating Whisper Service...")
        self.client =  Minio(
            "localhost:9000",
            access_key=minio_user,
            secret_key=minio_password,
            secure=False
        )
        self.bucket_name = "bucket"

    def get(self, object_name, file_path):
        self.client.fget_object(self.bucket_name, object_name, file_path)

    def save(self, object_name, file_path):
        self.client.fput_object(self.bucket_name, object_name, file_path)
