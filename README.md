# Thing Manager
[![Drone (cloud)](https://img.shields.io/drone/build/I1820/tm.svg?style=flat-square)](https://cloud.drone.io/I1820/tm)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/tm)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/i1820/tm.svg?style=flat-square)](https://hub.docker.com/r/i1820/tm)

## Introduction
Thing manager manages I1820 Things and their properties.
Things belong to the projects, but this component doesn't validate this relationship so other services
must verify project identification and existence before calls this project APIs.
