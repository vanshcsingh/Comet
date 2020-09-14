import tensorflow as tf
import numpy as np
from functools import reduce
from concurrent import futures
import logging
import sys
from termcolor import colored

import grpc
import service_pb2_grpc as pb_grpc
import service_pb2 as pb

old_model = tf.keras.models.load_model('./checkpoints/seq_model')

def convertToLabels(predictResults):
    results_int = np.argmax(predictResults, axis=1)
    return list(map(str, results_int))

class Service(pb_grpc.ServiceServicer):
    def Predict(self, request, context):
        input_batch = np.stack(list(map(lambda x: np.reshape(x.pixels, (28,28,1)), request.images)))
        prediction = old_model.predict(input_batch)
        labels = convertToLabels(prediction)

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
        print (colored("Please supply one argument: desired port", "red"))
        sys.exit(1)

    serve(sys.argv[1])