# Link
[![Travis branch](https://img.shields.io/travis/com/I1820/pm/master.svg?style=flat-square)](https://travis-ci.com/I1820/pm)
[![Go Report](https://goreportcard.com/badge/github.com/I1820/link?style=flat-square)](https://goreportcard.com/report/github.com/I1820/link)
[![Buffalo](https://img.shields.io/badge/powered%20by-buffalo-blue.svg?style=flat-square)](http://gobuffalo.io)
[![Maintainability](https://api.codeclimate.com/v1/badges/83743166f3b5ff30c429/maintainability)](https://codeclimate.com/github/I1820/link/maintainability)

## Introduction

Link component of I1820 platfrom. This service collects
raw data from bottom layer (protocols), stores them into mongo database
and decodes them using user selected decoder.
This service also sends data into bottom layer (protocols) after
encoding them using user selected encoder.
