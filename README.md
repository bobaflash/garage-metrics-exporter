# Garage Metric exporter

Since the garage metrics endpoint does not export metrics about the actual bucket size, this small program was born.
It provides prometheus metrics about the actual bucket size and should be easily extendable.

## Buliding

To build the binary run following command

```shell
go build -o garage-metrics-exporter cmd/main.go 
```

## Configuration

Following environment variables are available

| name                    | default               | description                                                             |
| ----------------------- | --------------------- | ----------------------------------------------------------------------- |
| GARAGE_BASE_URL         | http://localhost:3903 | The listen address of the garage admin endpoint                         |
| TOKEN                   |                       | The  token for the garage admin endpoint                                |
| LISTEN_SOCKET           | :3905                 | The socket this exporter listens on                                     |
| UPDATE_INTERVAL_SECONDS | 30                    | How often should the exporter request the garage admin API (in seconds) |
