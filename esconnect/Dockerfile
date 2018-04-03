# $ docker build -t oldenbur/esconnect .

FROM confluentinc/cp-kafka-connect

WORKDIR /esconnect
COPY . .

CMD ["/usr/bin/connect-standalone","connect-standalone.properties","connect-elasticsearch.properties"]

