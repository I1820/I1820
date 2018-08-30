# Data Manager
[![Travis branch](https://img.shields.io/travis/com/I1820/dm/master.svg?style=flat-square)](https://travis-ci.com/I1820/dm)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/dm?style=flat-square)](https://goreportcard.com/report/github.com/I1820/dm)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/2cda8cad3c7b46879da2544c1057c91f)](https://www.codacy.com/app/i1820/dm?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/dm&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/dm)

## Introduction
DM is a Data Manager component of the I1820 platform.
It has some useful built-in queries that can returns data from the database (MongoDB).
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
