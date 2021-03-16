from botocore.client import BaseClient
import boto3
import os
import json
import pickle

from botocore.exceptions import ClientError

access_key = 'AWS_ACCESS_KEY_ID'
access_value = os.getenv(access_key)
secret_key = 'AWS_SECRET_ACCESS_KEY'
secret_value = os.getenv(secret_key)

bucket_name = "ibmqo"


class S3client:

    # uses credentials from environment

    def s3_connect(self) -> BaseClient:
        s3 = boto3.client(
            's3',
            aws_access_key_id=access_value,
            aws_secret_access_key=secret_value,
        )
        return s3

    def put(self, bucket_name, key, obj):
        print("Storing")
        s3 = boto3.client('s3')
        # Serialize the object
        serializedListObject = json.dumps(obj)
        s3.put_object(Bucket=bucket_name, Key=key,
                      Body=serializedListObject)

    def get(self, bucket_name, key):
        s3 = self.s3_connect()
        try:
            response = s3.get_object(Bucket=bucket_name, Key=key)
            content = response['Body']
            jsonObject = json.loads(content.read())
            print(jsonObject)
            return jsonObject
        except ClientError as ex:
            if ex.response['Error']['Code'] == 'NoSuchKey':
                print('No object found - returning empty')
            else:
                raise
