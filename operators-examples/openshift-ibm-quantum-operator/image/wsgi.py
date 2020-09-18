from flask import Flask
from controller import Controller
from flask import request
import json
import pickle
import logging
from circuit_wrapper import CircuitWrapper


application = Flask(__name__)
controller = Controller()


@application.route('/submitJob', methods=['POST'])
def api_submitJob():
    logging.info("In sever request received: %s %s " %
                 (request, request.content_length))

    j_request = request.get_data()
    response = controller.handleJobRequest(j_request)
    return response


@application.route('/status', methods=['POST'])
def api_getStatus():
    logging.info("In sever request received: %s %s " %
                 (request, request.content_length))
    j_request = request.get_data()
    response = controller.handleStatusRequest(j_request)
    return response


@application.route('/', methods=['GET'])
def api_get():
    logging.info("In sever request received: %s %s " %
                 (request, request.content_length))

    return "Hi"


if __name__ == '__main__':
    logging.getLogger().setLevel(logging.INFO)
    application.run()
