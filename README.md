# Data Manager
[![Drone (cloud)](https://img.shields.io/drone/build/I1820/dm.svg?style=flat-square)](https://cloud.drone.io/I1820/dm)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/dm)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/i1820/dm.svg?style=flat-square)](https://hub.docker.com/r/i1820/dm)

## Introduction
DM is a Data Manager component of the I1820 platform. It handles data that are coming from RabbitMQ and stores them.
It also has some useful built-in queries that can returns data from the database (MongoDB).
We plan to create [Grafana](https://grafana.com/) plugin for it.

## Profiler
Enable MongoDB buit-in profiler:

```
use i1820
db.setProfileLevel(2)
```

The profiling results will be in a special capped collection called `system.profile`
which is located in the database where you executed the `setProfileLevel` command.

```
db.system.profile.find().pretty()
```
## Up and Running
To build this module from source do the following steps

1. Make sure MongoDB is up and running.

2. Setup MongoDB using the scripts provided in `mongodb/`.

3. Run :runner:
```sh
go build
./dm
```
