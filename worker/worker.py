import json

import grpc
import asyncio
import worker_pb2
import worker_pb2_grpc
from sybil import Sybil

model = Sybil()


class WorkerServicer(worker_pb2_grpc.WorkerServicer):
    async def Inference(self, request, context):
        res = {
            "prediction": [0.35] * 6,
            "threshold": 0.5
        }
        return worker_pb2.InferenceResponse(series=request.series, result=json.dumps(res))


async def serve():
    server = grpc.aio.server()
    worker_pb2_grpc.add_WorkerServicer_to_server(WorkerServicer(), server)
    server.add_insecure_port('[::]:50051')
    await server.start()
    try:
        await server.wait_for_termination()
    finally:
        await server.stop(3)


if __name__ == '__main__':
    asyncio.run(serve())
