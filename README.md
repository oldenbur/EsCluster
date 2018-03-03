# EsCluster
Elasticsearch cluster of docker containers fed by kafka

Docker Components
* [elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)
* [kafka](https://github.com/wurstmeister/kafka-docker)

```
# clone this repo
$ git clone [...]

# get and run elasticsearch
$ docker pull docker.elastic.co/elasticsearch/elasticsearch-oss:6.1.1
$ docker-compose up

# set up and run kafka - see http://wurstmeister.github.io/kafka-docker/
$ mkdir wurstmeister; cd !$
$ git clone https://github.com/wurstmeister/kafka-docker.git
$ cd kafka-docker
$ vim docker-compose.yml # modify KAFKA_ADVERTISED_HOST_NAME and zookeeper to mac IP
$ docker-compose up -d

$ ./start-kafka-shell.sh 192.168.99.36 192.168.99.36:2181
# $KAFKA_HOME/bin/kafka-topics.sh --create --topic topic1 --partitions 1 --zookeeper $ZK --replication-factor 1
# $KAFKA_HOME/bin/kafka-topics.sh --list --zookeeper $ZK
# $KAFKA_HOME/bin/kafka-console-producer.sh --topic topic1 --broker-list=`broker-list.sh`


# create a docker swarm vm:
$ docker-machine create --driver virtualbox --virtualbox-hostonly-cidr "10.10.10.1/24" myvm1

```

