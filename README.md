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

## Installation
Replace the Token in the service file with your personal token.
Usually, this token is configured in the file `/etc/garage.toml`

Maybe you also need to change the environment variables in the service to your needs.

```shell
# Download file from repo
wget https://github.com/bobaflash/garage-metrics-exporter/releases/download/v0.0.1/garage-metrics-exporter-linux-amd64.tar.gz

# Extract file and copy to folder /usr/local/bin/
tar -xvzf garage-metrics-exporter-linux-amd64.tar.gz
cp garage-metrics-exporter /usr/local/bin/

# Create systemd service
cat <<-EOF > /etc/systemd/system/garage-metrics-exporter.service
[Unit]
Description=Garage Metrics Exporter
After=network-online.target
Wants=network-online.target

[Service]
Environment='TOKEN=<<TOP SECRET TOKEN>>' 'UPDATE_INTERVAL_SECONDS=60'
ExecStart=/usr/local/bin/garage-metrics-exporter

[Install]
WantedBy=multi-user.target
EOF


# Enable and start service
systemctl daemon-reload
systemctl enable garage-metrics-exporter.service
systemctl start garage-metrics-exporter.service

# Check status
systemctl status garage-metrics-exporter.service
```