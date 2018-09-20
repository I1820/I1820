# Project Manager
[![Travis branch](https://img.shields.io/travis/com/I1820/pm/master.svg?style=flat-square)](https://travis-ci.com/I1820/pm)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/7dd562018dbc45f4a069c12c48195add)](https://www.codacy.com/app/i1820/pm?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/pm&amp;utm_campaign=Badge_Grade)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/pm?style=flat-square)](https://goreportcard.com/report/github.com/I1820/pm)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/pm)


## Introduction
PM is a project manager component of the I1820 platform.
It builds things, projects, things to project relationship, and project's dockers.
Each project consists of two dockers one of them provides a sandbox for user scripts that are in python and based on [ElRunner](https://github.com/I1820/ElRunner) and another
runs redis as in-memory storage for user scripts.
It builds these dockers in localhost and uses Linux sockets for communicating with docker host.

PMs can run on many hosts to provide load balancing. To distribute requests among them, you can use [vulcand](https://vulcand.readthedocs.io/en/latest/quickstart.html#quick-start). this feature still in development phase so it would be better not to use it now :joy:

This component provides API based on HTTP ReST so other components can utilize these APIs for creating and destroying things and projects.

PM requires only MongoDB to persist things and projects data.

## Assets
An asset is a new concept that is added recently to PM. assets are sensors or actuators that are connected into things.
Assets can send or receive data based on their kind.

## Up and Running

To use this module you can use its docker or build from source
after that, you must do the following things to provide the foundation for project creation.

1. Non-Root user access for docker
```sh
sudo usermod -aG docker $USER
```

2. Create projects network
```sh
docker network create -d bridge --subnet 192.168.72.0/24 --gateway 192.168.72.1 i1820
```

3. Pull required images
```sh
docker pull i1820/elrunner
docker pull redis:alpine
```
