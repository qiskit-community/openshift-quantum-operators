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


# if __name__ == "__main__":
#     s3_connect()
#     bucket_name = "ibmqo"
#     key = "cd5cc452-c2e0-42cb-8f37-fd381e3aa472"
#     obj = {
#         "qobj_id": "f3ccd747-8812-41fb-974e-cddc447c63f4",
#         "header": {},
#         "config": {
#             "shots": 1024,
#             "memory": False,
#             "parameter_binds": [],
#             "init_qubits": True,
#             "memory_slots": 2,
#             "n_qubits": 2
#         },
#         "schema_version": "1.3.0",
#         "type": "QASM",
#         "experiments": [
#             {
#                 "config": {
#                     "n_qubits": 2,
#                     "memory_slots": 2
#                 },
#                 "header": {
#                     "qubit_labels": [
#                         [
#                             "q",
#                             0
#                         ],
#                         [
#                             "q",
#                             1
#                         ]
#                     ],
#                     "n_qubits": 2,
#                     "qreg_sizes": [
#                         [
#                             "q",
#                             2
#                         ]
#                     ],
#                     "clbit_labels": [
#                         [
#                             "c",
#                             0
#                         ],
#                         [
#                             "c",
#                             1
#                         ]
#                     ],
#                     "memory_slots": 2,
#                     "creg_sizes": [
#                         [
#                             "c",
#                             2
#                         ]
#                     ],
#                     "name": "circuit7",
#                     "global_phase": 0
#                 },
#                 "instructions": [
#                     {
#                         "name": "u2",
#                         "params": [
#                             0,
#                             3.141592653589793
#                         ],
#                         "qubits": [
#                             0
#                         ]
#                     },
#                     {
#                         "name": "cx",
#                         "qubits": [
#                             0,
#                             1
#                         ]
#                     },
#                     {
#                         "name": "measure",
#                         "qubits": [
#                             0
#                         ],
#                         "memory": [
#                             0
#                         ]
#                     },
#                     {
#                         "name": "measure",
#                         "qubits": [
#                             1
#                         ],
#                         "memory": [
#                             1
#                         ]
#                     }
#                 ]
#             }
#         ]
#     }
#     store(bucket_name, key, obj)
#     get(bucket_name, key)
