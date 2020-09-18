from qiskit import *
from qiskit.visualization import plot_histogram
from qiskit.tools.monitor import job_monitor
from scheduler import Scheduler
import requests
import json
from circuit_wrapper import CircuitWrapper
import pickle
import logging
import os
import configparser
from qiskit.providers.ibmq.exceptions import *


class Controller:

    def __init__(self):

        token = ''

        # config = self.loadSecret()
        # IBMQ.save_account(config.get('token'), overwrite=True)
        IBMQ.save_account(token, overwrite=True)
        provider = IBMQ.load_account()
        self.scheduler = Scheduler(provider)

    def loadSecret():
        print('Loading secret')
        config = configparser.ConfigParser()

        try:
            f = open("/tmp/secrets/qiskitsecret/qiskit-secret.cfg", 'r')
            config.readfp(f)
        finally:
            f.close()

        options = {
            'auth_version':         3,
            'token':          config['AUTH TOKENS']['token'],
        }
        return options

    def handleJobRequest(self, j_req):
        logging.info("handleJobRequest in Controller")

        circuitWrapper = pickle.loads(j_req)

        circuit = circuitWrapper.circuit
        circuitName = circuitWrapper.circuitName
        qubits = circuitWrapper.qubits

        q_job = self.scheduler.schedule(circuit, qubits, circuitName)
        circuitWrapper.backend = q_job.backend()

        circuitWrapper = self.getStatus(q_job, circuitWrapper)
        logging.info("Job %s is submitted with status %s" % (circuitWrapper.jobId,
                                                             circuitWrapper.status))
        data = pickle.dumps(circuitWrapper)
        return data

    def getStatus(self, job, circuitWrapper):
        # Get the status from the IBM backend now and add an object store to fetch results
        status = job.status()
        circuitWrapper.jobId = job.job_id()

        if job.done():
            circuitWrapper.answer = job.result().get_counts(circuitWrapper.circuit)
        else:
            job.refresh()

        circuitWrapper.status = job.status().value
        return circuitWrapper

    def handleStatusRequest(self, j_req):
        logging.info("handleStatusRequest in Controller")

        circuitWrapper = pickle.loads(j_req)
        jobId = circuitWrapper.jobId

        provider = IBMQ.get_provider()
        logging.debug("Provider: %s" % provider)

        backend = provider.get_backend(circuitWrapper.backend.name())
        logging.debug("Backend: %s" % backend)

        retrieved_job = backend.retrieve_job(jobId)
        logging.debug('Job retrieved using jobId %s' % jobId)

        circuitWrapper = self.getStatus(retrieved_job, circuitWrapper)
        logging.info("Job %s is submitted with status %s" % (circuitWrapper.jobId,
                                                             circuitWrapper.status))
        data = pickle.dumps(circuitWrapper)

        return data
