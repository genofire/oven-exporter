package main

import (
	"flag"
	"encoding/base64"

	"dev.sum7.eu/genofire/golang-lib/file"
	"github.com/bdlm/log"
)

type configData struct {
	Token string `toml:"token"`
	Host  string `toml:"host"`
}

func main() {
	configPath := "config.toml"

	flag.StringVar(&configPath, "c", configPath, "path to configuration file")

	flag.Parse()

	config := &configData{}
	if err := file.ReadTOML(configPath, config); err != nil {
		log.Panicf("open config file: %s", err)
	}
	config.Token = base64.StdEncoding.EncodeToString([]byte(config.Token))
	config.fetch()
}
