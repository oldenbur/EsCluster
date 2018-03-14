## EsCluster
Dockerized kafka elasticsearch streaming and persistence pipeline.
* confluent kafka [docker stack](https://docs.confluent.io/current/installation/docker/docs/quickstart.html)
* golang elasticsearch json connector [kclient](#kclient-docker-image)

#### Install kafka Command Line Tools
See https://kafka.apache.org/quickstart
* Mac: `wget http://apache.claz.org/kafka/1.0.1/kafka_2.11-1.0.1.tgz`

#### kclient Docker Image
Based on the golang docker development [guide](https://hub.docker.com/_/golang/).
```
$ cd kclient
$ docker build -t oldenbur/kclient .
```

##### Local development
To develop kclient locally, install the rdkafka library:
* Mac: `brew install librdkafka`

#### Run the Stack
Run a simple end-to-end integration test where the kafka console producer run on the
host machine feeds integers through a kafka topic to kclient:
```
$ cd esCluster
$ export DOCKER_MACHINE_IP=[mac_ip]
$ docker-compose up -d
$ cd kafka_2.11-1.0.1
$ bin/kafka-topics.sh --create --topic foo --partitions 1 --replication-factor 1 --if-not-exists --zookeeper ${DOCKER_MACHINE_IP}:32181
$ seq 42 | bin/kafka-console-producer.sh --request-required-acks 1 --broker-list ${DOCKER_MACHINE_IP}:29092 --topic foo && echo 'Produced 42 messages.'
```
