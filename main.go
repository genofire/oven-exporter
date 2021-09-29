package main

import (
	"flag"
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"dev.sum7.eu/genofire/oven-exporter/api"
)

type configData struct {
	log    *zap.Logger
	Log    *zap.Config `toml:"log"`
	API    *api.Client `toml:"api"`
	Listen string      `toml:"listen"`
}

func main() {
	configPath := "config.toml"

	log := zap.L()

	flag.StringVar(&configPath, "c", configPath, "path to configuration file")

	flag.Parse()

	config := &configData{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panic("open config file", zap.Error(err))
	}
	if config.Log != nil {
		l, err := config.Log.Build()
		if err != nil {
			log.Panic("generate logger from config", zap.Error(err))
		}
		log = l
	}
	config.log = log
	//config.SetLogger(log)
	config.API.SetToken(config.API.Token)

	prometheus.MustRegister(config)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		log.Fatal("crash webserver", zap.Error(err))
	}
}
