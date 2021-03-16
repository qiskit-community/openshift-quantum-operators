from fastapi import FastAPI, Request, status
from fastapi.encoders import jsonable_encoder
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from controller import Controller

import pickle
import json

import numpy

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.post("/submit/")
async def submit(request: Request):
    print("request: ", request)
    body = await request.json()
    print("--------------------------------------------------------------")
    print("Body: \n", body)
    print("Type of body", type(body))
    ctl = Controller()
    # ctl.submit(body)

    return ctl.submit(body)


@app.post("/status/")
async def status(request: Request):
    print("request: ", request)
    body = await request.json()
    print("--------------------------------------------------------------")
    print("Body: \n", body)
    print("Type of body", type(body))
    ctl = Controller()

    return ctl.status(body)
