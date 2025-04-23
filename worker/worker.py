import json
import logging
import os
import shutil
import tempfile
import zipfile
import grpc
import asyncio
import requests
import worker_pb2
import worker_pb2_grpc

from pydicom.misc import is_dicom
from sybil import Sybil, Serie
from sybil.utils.device_utils import get_default_device

model = Sybil(device=get_default_device(), cache='./models')


def process(path: str) -> list[float] | None:
    tmp_dir = tempfile.mkdtemp()
    zip_path = os.path.join(tmp_dir, 'archive.zip')
    extract_dir = os.path.join(tmp_dir, 'extracted')
    try:
        with requests.get(path, stream=True) as r:
            r.raise_for_status()
            with open(zip_path, 'wb') as f:
                for chunk in r.iter_content(chunk_size=8192):
                    f.write(chunk)
        with zipfile.ZipFile(zip_path, 'r') as zip_ref:
            zip_ref.extractall(extract_dir)
            paths = []
            for f in os.listdir(extract_dir):
                file_path = os.path.join(extract_dir, f)
                if is_dicom(file_path):
                    paths.append(file_path)

            paths.sort()
            predictions = model.predict(series=Serie(paths))
            return predictions.scores[0]
    finally:
        shutil.rmtree(tmp_dir)


class WorkerServicer(worker_pb2_grpc.WorkerServicer):
    async def Inference(self, request, context):
        if request.cpu:
            model.to(device='cpu')
        else:
            model.to(device=get_default_device())

        prediction = process(request.path)
        res = {
            "prediction": prediction,
            "threshold": 0.5
        }
        return worker_pb2.InferenceResponse(series=request.series, result=json.dumps(res))


async def serve():
    server = grpc.aio.server()
    worker_pb2_grpc.add_WorkerServicer_to_server(WorkerServicer(), server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    print('Server started')
    try:
        await server.wait_for_termination()
    finally:
        await server.stop(3)


if __name__ == '__main__':
    os.environ.setdefault("PYTORCH_ENABLE_MPS_FALLBACK", "1")
    asyncio.run(serve())
