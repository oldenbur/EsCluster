## EsCluster
Elasticsearch cluster of docker containers fed by kafka.
* confluent kafka [docker stack](https://docs.confluent.io/current/installation/docker/docs/quickstart.html)
* golang json connector [kclient](#kclient)

### kclient
Golang implementation of a kafka consumer.

#### Install librdkafka
* Mac: `brew install librdkafka`

#### Install kafka
See https://kafka.apache.org/quickstart
* Mac: `wget http://apache.claz.org/kafka/1.0.1/kafka_2.11-1.0.1.tgz`

#### Docker
Based on the golang docker development [guide](https://hub.docker.com/_/golang/).
```
$ cd kclient
$ docker build -t oldenbur/kclient .
$ docker run -it --rm --name kclient oldenbur/kclient
```

