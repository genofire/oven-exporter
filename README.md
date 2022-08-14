# Oven-Exporter
An Prometheus Exporter for OvenMediaEngine

(it provides also a small API-Client for OvenMediaEngine)
Be welcome to improve it.

## Configure OvenMediaEngine

This Exporter use the REST-API of OvenMediaEngine,
to setting it up, that a look in there Documentation [OvenMediaEngine REST-API](
https://airensoft.gitbook.io/ovenmediaengine/rest-api).

## Setup Exporter

### Compile

Install [golang](https://golang.org/doc/install).

Run:
`go install -v dev.sum7.eu/genofire/oven-exporter@latest`

### Configuration
Read comments in [config_example.toml](config_example.toml) for more information.

Maybe a good place to store this file is: `/etc/ovenmediaengine/exporter.conf`

OR use env variables:
```
OVEN_E_LISTEN=:8080
OVEN_E_API__URL=http://1.2.3.4:8081
OVEN_E_API__TOKEN=ome-access-token
OVEN_E_API__DEFAULT_VHOST=
OVEN_E_API__DEFAULT_APP=
```

(File read could be disabled by call `oven-exporter -c ''`

### Startup
Create a systemd.service file e.g. under `/etc/systemd/system/oven-exporter.service` with maybe a content like this:

```ini
[Unit]
Description = Prometheus exporter for OvenMediaEngine

[Service]
Type=simple
ExecStart=/usr/local/bin/oven-exporter -c /etc/ovenmediaengine/exporter.conf
Restart=always
RestartSec=5s
Environment=PATH=/usr/bin:/usr/local/bin

[Install]
WantedBy=multi-user.target
```
PS: maybe you need to adjust the binary path and configuration path.

Start and enable on boot:
`systemctl enable --now oven-exporter.service`
