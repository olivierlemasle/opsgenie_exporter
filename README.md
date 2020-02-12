# Opsgenie exporter for Prometheus

[![Docker Pulls](https://img.shields.io/docker/pulls/olem/opsgenie_exporter.svg?maxAge=604800)](https://hub.docker.com/r/olem/opsgenie_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/olivierlemasle/opsgenie_exporter)](https://goreportcard.com/report/github.com/olivierlemasle/opsgenie_exporter)

## Running using Docker

```
docker run -d -p 3000:3000 -v /path/on/host/conf/lamp.conf:/conf/lamp.conf olem/opsgenie_exporter
```
