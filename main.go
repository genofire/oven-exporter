package main

import (
	"flag"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"codeberg.org/Mediathek/oven-exporter/api"
)

var configExtParser = map[string]koanf.Parser{
	".json": json.Parser(),
	".toml": toml.Parser(),
	".yaml": yaml.Parser(),
	".yml":  yaml.Parser(),
}

type configData struct {
	log    *zap.Logger
	Log    *zap.Config `config:"log"`
	API    api.Client  `config:"api"`
	Listen string      `config:"listen"`
}

func main() {

	configPath := "config.toml"

	log, _ := zap.NewProduction()

	flag.StringVar(&configPath, "c", configPath, "path to configuration file")
	flag.Parse()

	k := koanf.New("/")

	if configPath != "" {
		fileExt := filepath.Ext(configPath)
		parser, ok := configExtParser[fileExt]
		if !ok {
			log.Panic("unsupported file extension:",
				zap.String("config-path", configPath),
				zap.String("file-ext", fileExt),
			)
		}
		if err := k.Load(file.Provider(configPath), parser); err != nil {
			log.Panic("load file config:", zap.Error(err))
		}
	}

	if err := k.Load(env.Provider("OVEN_E_", "/", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "OVEN_E_")), "__", "/", -1)
	}), nil); err != nil {
		log.Panic("load env:", zap.Error(err))
	}

	config := &configData{}
	if err := k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: "config"}); err != nil {
		log.Panic("reading config", zap.Error(err))
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
