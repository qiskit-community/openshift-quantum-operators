
from qiskit import BasicAer
from qiskit import QuantumCircuit
from qiskit import transpile, assemble
from qiskit.qobj.qasm_qobj import QasmQobj as QasmQobj

import pickle
import json
import numpy

from s3client.connect import S3client


class Controller:

    def status(self, qasm_dict={}):
        print("Recieved: ", qasm_dict)
        job_id = qasm_dict['qobj_id']
        s3client = S3client()
        bucket_name = "ibmqo"
        obj = s3client.get(bucket_name, job_id)
        return obj

    def submit(self, qasm_dict={}):
        print("Recieved qasm_dict: ", qasm_dict)
        qasm_ojb = QasmQobj.from_dict(qasm_dict)
        backend_sim = BasicAer.get_backend('qasm_simulator')

        result = backend_sim.run(qasm_ojb).result()
        print("Result: ", result)

        result_dict = result.to_dict()
        print("result_dict: ", result_dict)

        result_json = json.dumps(result_dict, cls=QobjEncoder)
        print("result_json \n: ", result_json)

        bucket_name = "ibmqo"
        key = "cd5cc452-c2e0-42cb-8f37-fd381e3aa472"
        obj = result_json
        s3client = S3client()
        s3client.put(bucket_name, key, obj)
        return result_json


class QobjEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, numpy.int32):
            return obj.item()
        if isinstance(obj, numpy.ndarray):
            return obj.tolist()
        if isinstance(obj, complex):
            return (obj.real, obj.imag)
        return json.JSONEncoder.default(self, obj)
