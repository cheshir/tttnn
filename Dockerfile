FROM tensorflow/tensorflow

RUN pip install --upgrade pip h5py keras

RUN add-apt-repository -y ppa:longsleep/golang-backports; \
    apt update; \
    apt install -y vim wget golang-go;

RUN mkdir -p /go/src
ENV GOPATH=/go

# CPU only version.
RUN wget -O libtensorflow.tar.gz https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-1.12.0.tar.gz; \
    tar -xzf libtensorflow.tar.gz -C /usr/local

ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

EXPOSE 8536
