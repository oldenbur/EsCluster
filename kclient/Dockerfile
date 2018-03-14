# $ docker build -t oldenbur/kclient .

FROM golang:latest

# Install librdkafka
RUN mkdir -p /root/src && cd /root/src
WORKDIR /root/src
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/src/librdkafka
RUN ./configure
RUN make
RUN make install

ENV LD_LIBRARY_PATH /usr/local/lib

WORKDIR /go/src/kclient
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["kclient"]
