# EsCluster
Elasticsearch cluster of docker containers fed by kafka

## kclient

### Install librdkafka
* Mac: `brew install librdkafka`

### Docker
Based on the golang docker development [guide](https://hub.docker.com/_/golang/).
```
$ cd kclient
$ docker build -t oldenbur/kclient .
$ docker run -it --rm --name kclient oldenbur/kclient
```

