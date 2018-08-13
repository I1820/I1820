# Project Manager
[![Travis branch](https://img.shields.io/travis/com/I1820/pm/master.svg?style=flat-square)](https://travis-ci.com/I1820/pm)
[![Maintainability](https://api.codeclimate.com/v1/badges/e8583a735941b7d9a505/maintainability)](https://codeclimate.com/github/I1820/pm/maintainability)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/pm?style=flat-square)](https://goreportcard.com/report/github.com/I1820/pm)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)


## Introduction
PM is a project manager component of the I1820 platform.
it builds projects and their dockers. docker provides a sandbox for user scripts that are in python.

## Installation

1. Non-Root user access for docker
```sh
sudo usermod -aG docker $USER
```

2. Create ISRC network
```sh
docker network create i1820
```

3. Pull required images
```sh
docker pull i1820/elrunner
docker pull redis:alpine
```
