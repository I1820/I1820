# Link
[![Travis branch](https://img.shields.io/travis/com/I1820/link/master.svg?style=flat-square)](https://travis-ci.com/I1820/link)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/link?style=flat-square)](https://goreportcard.com/report/github.com/I1820/link)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1bdf3a4f0b294e9e92f15211ba894ef4)](https://www.codacy.com/app/i1820/link?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=I1820/link&amp;utm_campaign=Badge_Grade)

## Introduction

Link component of I1820 platfrom. This service collects
raw data from bottom layer (protocols), stores them into mongo database
and decodes them using user selected decoder.
This service also sends data into bottom layer (protocols) after
encoding them using user selected encoder.

Link uses MQTT for communicating with the bottom layer and this communication can be customized
using Protocol interface which defined in `app/app.go`.
