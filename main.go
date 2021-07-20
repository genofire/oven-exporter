package main

import (
	"flag"
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"dev.sum7.eu/genofire/oven-exporter/api"
)

type configData struct {
	API    *api.Client `toml:"api"`
	Listen string      `toml:"listen"`
}

func main() {
	configPath := "config.toml"

	flag.StringVar(&configPath, "c", configPath, "path to configuration file")

	flag.Parse()

	config := &configData{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panicf("open config file: %s", err)
	}
	config.API.SetToken(config.API.Token)

	prometheus.MustRegister(config)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(config.Listen, nil))
}
