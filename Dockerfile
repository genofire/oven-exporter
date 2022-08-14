##
# Compile application
##
FROM golang:alpine AS build-env
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
COPY --from=build-env ["/etc/ssl/cert.pem", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=build-env /app/oven-exporter /oven-exporter
WORKDIR /
ENTRYPOINT ["/oven-exporter", "-c", ""]
