FROM ubuntu:14.04
RUN apt-get update \
    && apt-get install -y curl\
    && curl -O https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz \
    && tar -xvf go1.8.1.linux-amd64.tar.gz \
    && rm go1.8.1.linux-amd64.tar.gz \
    && mv go /usr/local \
    && mkdir /root/MapReduce-With-Amazon_S3
WORKDIR /root
ENV GOPATH /root/MapReduce-With-Amazon_S3
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV AWS_REGION ap-northeast-1
ENV AWS_ACCESS_KEY_ID sssss
ENV AWS_SECRET_ACCESS_KEY sssss
COPY src /root/MapReduce-With-Amazon_S3/src
WORKDIR /root/MapReduce-With-Amazon_S3
CMD ["go", "run", "src/main/primality_mr.go", "worker", "172.31.31.161:7777", "0.0.0.0:7778"]
