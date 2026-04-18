import os

minio_user = os.environ.get('MINIO_ROOT_USER', 'admin')
minio_password = os.environ.get('MINIO_ROOT_PASSWORD', "adminpass")
hf_token = os.environ.get('HF_TOKEN')

print(f"User: {minio_user}")