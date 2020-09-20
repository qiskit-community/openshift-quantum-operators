from qiskit.tools.monitor import job_monitor
from qiskit import IBMQ, Aer, assemble, transpile
from qiskit.providers.ibmq import least_busy
from qiskit import QuantumCircuit, ClassicalRegister, QuantumRegister, execute

import logging


class Scheduler:
    provider = None

    def __init__(self, provider):
        self.provider = provider

    def schedule(self, circuit, qubits, circuitName):

        logging.info('Scheduling job to IBM Q in progess')

        backend = least_busy(self.provider.backends(filters=lambda x: x.configuration().n_qubits >= qubits and
                                                    not x.configuration().simulator and x.status().operational == True))

        # backend = self.provider.backends.ibmq_qasm_simulator
        logging.debug("Least busy backend: %s" % backend)
        logging.debug("In scheduler number of qubits requested: %s" % qubits)
        logging.debug("In scheduler: ", circuit)

        qobj = assemble(transpile(circuit, backend=backend), backend=backend)
        job = backend.run(qobj)

        logging.debug("Job is submitted with jobId %s" % job.job_id())

        return job
