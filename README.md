# ISRC Uplink
## Introduction
Uplink service of ISRC platfrom. This service collects
raw data from bottom layer, stores them into mongo database
and decodes them using decoder on runner platform.

## Running
MongoDB

```sh
docker run -ti --rm -p 27017:27017 mongo
```

Vernemq

```sh
docker run -e "DOCKER_VERNEMQ_ALLOW_ANONYMOUS=on" --rm -ti -p 1883:1883 --name vernemq1 erlio/docker-vernemq
```
