FROM tensorflow/tensorflow

RUN apt update && apt install -y vim
RUN pip install --upgrade pip h5py keras
