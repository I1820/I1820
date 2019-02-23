# Thing Manager
[![Travis branch](https://img.shields.io/travis/com/I1820/tm/master.svg?style=flat-square)](https://travis-ci.com/I1820/tm)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/tm?style=flat-square)](https://goreportcard.com/report/github.com/I1820/tm)
[![Echo](https://img.shields.io/badge/powered%20by-echo-blue.svg?style=flat-square)](https://echo.labstack.com/)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/I1820/tm)

## Introduction
Thing manager manages I1820 Things and their properties. Each thing in I1820 have the following
properties:

- Assets
- Connectivities
- Tags
- etc.

Things belong to the projects, but this component doesn't validate this relationship so other services
must verify project identification and existence before calls this project APIs.
