package main

import (
	"flag"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"

	"dev.sum7.eu/genofire/oven-exporter/api"
)

type configData struct {
	API *api.Client `toml:"api"`
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
	fetch(config.API)
}
