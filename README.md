[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/i1820/i1820.svg?style=flat-square)](https://hub.docker.com/r/i1820/i1820)
[![Drone (cloud)](https://img.shields.io/drone/build/I1820/I1820.svg?style=flat-square)](https://cloud.drone.io/I1820/I1820)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/I1820)


## Link
Link component of I1820 platfrom. This service collects
raw data from bottom layer (protocols), stores them into mongo database
and decodes them using user's selected decoder.
This service also sends data into bottom layer (protocols) after
encoding them using user's selected encoder.

Link uses MQTT for communicating with the bottom layer and this communication can be customized
using Protocol's interface which is defined in `protocols/protocol.go`.
