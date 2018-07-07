# Project Manager
[![Travis branch](https://img.shields.io/travis/aiotrc/pm/master.svg?style=flat-square)](https://travis-ci.org/aiotrc/pm)
[![Codacy grade](https://img.shields.io/codacy/grade/f536424b14cc4df5998f4ca0b356b661.svg?style=flat-square)](https://www.codacy.com/app/1995parham/pm?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aiotrc/pm&amp;utm_campaign=Badge_Grade)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)


## Introduction
PM is project manager for ISRC platform. It creates projects and corresponding runners (containers) for users.

## Installation
1. Non-Root user access for docker
```sh
sudo usermod -aG docker $USER
```
2. Create ISRC network
```sh
docker network create isrc
```
