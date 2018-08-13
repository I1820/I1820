# Data Manager
[![Travis branch](https://img.shields.io/travis/I1820/dm/master.svg?style=flat-square)](https://travis-ci.org/I1820/dm)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/dm?style=flat-square)](https://goreportcard.com/report/github.com/I1820/dm)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Maintainability](https://api.codeclimate.com/v1/badges/cf30f8b1aa2317c8b44e/maintainability)](https://codeclimate.com/github/I1820/dm/maintainability)

## Introduction
DM queries and returns data from database (mongodb).
it has grafana plugin for better data management.

## Profiler
Enable MongoDB buit-in profiler:

```
use isrc
db.setProfileLevel(2)
```

The profiling results in a special capped collection called `system.profile`
which is located in the database where you executed the `setProfileLevel` command.

```
db.system.profile.find().pretty()
```
