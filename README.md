# I1820
[![Drone (cloud)](https://img.shields.io/drone/build/I1820/I1820.svg?style=flat-square)](https://cloud.drone.io/I1820/I1820)

## Introduction

## Link
Link component of I1820 platfrom. This service collects
raw data from bottom layer (protocols), stores them into mongo database
and decodes them using user's selected decoder.
This service also sends data into bottom layer (protocols) after
encoding them using user's selected encoder.

Link uses MQTT for communicating with the bottom layer and this communication can be customized
using Protocol's interface which is defined in `protocols/protocol.go`.

## Thing Manager
Thing manager manages I1820 Things and their properties.
Things belong to the projects, but this component doesn't validate this relationship so other services
must verify project identification and existence before calls this project APIs.

## Data Manager
DM is a Data Manager component of the I1820 platform.
It has some useful built-in queries that can returns data from the database (MongoDB) to the API backend.

### Profiler
Enable MongoDB built-in profiler:

```
use i1820
db.setProfileLevel(2)
```

The profiling results will be in a special capped collection called `system.profile`
which is located in the database where you executed the `setProfileLevel` command.

```
db.system.profile.find().pretty()
```

## Project Manager

PM is a project manager component of the I1820 platform.
It builds things, projects, things to project relationship, and project's dockers.
Each project consists of two dockers one of them provides a sandbox for user scripts that are in python and based on [ElRunner](https://github.com/I1820/ElRunner) and another
runs redis as in-memory storage for user scripts.
It builds these dockers in localhost and uses Linux sockets for communicating with docker host.

PMs can run on many hosts to provide load balancing. To distribute requests among them, you can use [vulcand](https://vulcand.readthedocs.io/en/latest/quickstart.html#quick-start). this feature still in development phase so it would be better not to use it now :joy:

This component provides API based on HTTP ReST so other components can utilize these APIs for creating and destroying things and projects.

PM requires only MongoDB to persist things and projects data.

