FROM python:3.6-stretch
WORKDIR /home/vansh/z/Comet/src/comet/models/tensorflow/model
COPY . .
RUN pip install -r requirements.txt
ENTRYPOINT ["python3", "-u", "model.py"]