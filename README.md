# Sawzall

![Project](project.jpg)

Sawzall cuts through logs. It reads Nginx logs and summarizes them.

## Installation

Get the source and cd to it:

    go get -u github.com/vsimon/sawzall
    cd $(go env GOPATH)/src/github.com/vsimon/sawzall

## How to Build

There are two build dependencies that are required to build sawzall.

* Go 1.8+
* Make 3.81+

To build Sawzall:

    make

## How to Run

To run Sawzall:

    ./sawzall

By default, sawzall will read logs from `/var/log/nginx/access.log` and write
output to `/var/log/stats.log`.

A few optional parameters are supported:

    $ ./sawzall --help
    Usage of ./sawzall:
      -interval duration
            Interval to output summaries to output file (default 5s)
      -logFile string
            Nginx log file to read. (default "/var/log/nginx/access.log")
      -outputFile string
            Output file to write statsd-compatible summaries (default "/var/log/stats.log")

## How to Run Tests

To run the unit tests:

    make test

## Quick Demo

A demo is provided utilizing Docker Compose running a web server, load
generator, sawzall, and a simple log viewer.

    docker-compose up

Can also use Docker Compose to run the above instructions.

To build:

    docker-compose run app make

The resulting `sawzall` binary in the current directory will be a Linux binary.

To test:

    docker-compose run app make test


## High Level Overview

Sawzall is written in Golang. Sawzall is separated into a sawzall CLI program
located at `cmd/sawzall/main.go` and the sawzall library located at
`pkg/sawzall` which contains its software components and unit tests.

The Go dependency tool `dep` is used to manage and version a few dependencies.

At a high-level, Sawzall is composed of three main components:

* An `Ingester` which handles reading data from a file
* A `Summarizer` which handles summarizing the data from an Ingester
* A `Writer` which handles writing summaries from a Summarizer out to a file



