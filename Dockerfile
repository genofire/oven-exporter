##
# Compile application
##
FROM golang:latest AS build-env
WORKDIR /app
COPY . .
# ge dependencies
RUN go mod tidy
# build binary
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o oven-exporter


##
# Build Image
##
FROM scratch
COPY --from=build-env /app/oven-exporter /oven-exporter
COPY --from=build-env /app/config_example.toml /config.toml

WORKDIR /
ENTRYPOINT ["/oven-exporter"]
