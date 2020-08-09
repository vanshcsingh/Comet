import tensorflow as tf
import numpy as np
from functools import reduce
from concurrent import futures
import logging
import sys

import grpc
import service_pb2_grpc as pb_grpc
import service_pb2 as pb

old_model = tf.keras.models.load_model('./checkpoints/seq_model')

def convertToLabels(predictResults):
    results_int = np\
        .apply_along_axis(lambda x: np.where(x==1.)[0], axis=1, arr=predictResults)\
        .flatten()
    return list(map(str, results_int))

class Service(pb_grpc.ServiceServicer):
    def Predict(self, request, context):
        input_batch = np.stack(map(lambda x: x.vector, request.images))
        labels = convertToLabels(old_model.predict(input_batch))

        return pb.PredictReply(labels=labels)
        
def serve(port):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb_grpc.add_ServiceServicer_to_server(Service(), server)
    server.add_insecure_port('[::]:{}'.format(port))
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig()

    if len(sys.argv) != 2:
        print ("Please supply one argument: desired port")
        sys.exit(1)

    print ("arguments are", sys.argv)

    serve(sys.argv[1])