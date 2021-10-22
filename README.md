# Opsgenie exporter for Prometheus

[![Go Report Card](https://goreportcard.com/badge/github.com/olivierlemasle/opsgenie_exporter)](https://goreportcard.com/report/github.com/olivierlemasle/opsgenie_exporter)

## Flags

```
./opsgenie_exporter --help
```

- **`--web.listen-address`:** Address to listen on for web interface and telemetry. Defaults to `:3000`.
- **`--web.telemetry-path`:** Path under which to expose metrics. Defaults to `/metrics`.
- **`--web.disable-exporter-metrics`:** Exclude metrics about the exporter itself (`promhttp_*`, `process_*`, `go_*`).
- **`--log.level`**: Log level. Defaults to `info`.
- **`--log.format`**: Log format (`txt` or `json`). Defaults to `txt`.
- **`--config`:** Configuration file location. Defaults to `conf/lamp.conf`.
- **`--version`:** Display application version

## Configuration

`opsgenie_exporter` reads [Lamp configuration file](https://docs.opsgenie.com/docs/lamp-command-line-interface-for-opsgenie). Parameter `apiKey` is mandatory. An example can be found [here](./conf/lamp.conf.sample).

## Running using Docker

```
docker run -d \
  -p 3000:3000 \
  -v /path/on/host/conf/lamp.conf:/conf/lamp.conf \
  ghcr.io/olivierlemasle/opsgenie_exporter
```

Container images on ghcr.io registry: [opsgenie_exporter](https://github.com/olivierlemasle/opsgenie_exporter/pkgs/container/opsgenie_exporter)
