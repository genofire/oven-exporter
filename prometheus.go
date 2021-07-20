package main

import (
	"dev.sum7.eu/genofire/oven-exporter/api"
	"github.com/bdlm/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	promDescStatsVHostClientTotal   = prometheus.NewDesc("oven_vhost_client_total", "client count of current vhost", []string{"vhost"}, prometheus.Labels{})
	promDescStatsVHostClientMax     = prometheus.NewDesc("oven_vhost_client_max", "client max of current vhost", []string{"vhost"}, prometheus.Labels{})
	promDescStatsVHostBytesInTotal  = prometheus.NewDesc("oven_vhost_bytes_in_total", "total bytes in vhost", []string{"vhost"}, prometheus.Labels{})
	promDescStatsVHostBytesOutTotal = prometheus.NewDesc("oven_vhost_bytes_out_total", "total bytes out vhost", []string{"vhost"}, prometheus.Labels{})
	promDescStatsVHost              = []*prometheus.Desc{
		promDescStatsVHostClientTotal,
		promDescStatsVHostClientMax,
		promDescStatsVHostBytesInTotal,
		promDescStatsVHostBytesOutTotal,
	}

	promDescStatsAppClientTotal   = prometheus.NewDesc("oven_app_client_total", "client count of current app", []string{"vhost", "app"}, prometheus.Labels{})
	promDescStatsAppClientMax     = prometheus.NewDesc("oven_app_client_max", "client max of current app", []string{"vhost", "app"}, prometheus.Labels{})
	promDescStatsAppBytesInTotal  = prometheus.NewDesc("oven_app_bytes_in_total", "total bytes in app", []string{"vhost", "app"}, prometheus.Labels{})
	promDescStatsAppBytesOutTotal = prometheus.NewDesc("oven_app_bytes_out_total", "total bytes out app", []string{"vhost", "app"}, prometheus.Labels{})
	promDescStatsApp              = []*prometheus.Desc{
		promDescStatsAppClientTotal,
		promDescStatsAppClientMax,
		promDescStatsAppBytesInTotal,
		promDescStatsAppBytesOutTotal,
	}

	promDescStatsStreamClientTotal   = prometheus.NewDesc("oven_stream_client_total", "client count of current stream", []string{"vhost", "app", "stream"}, prometheus.Labels{})
	promDescStatsStreamClientMax     = prometheus.NewDesc("oven_stream_client_max", "client max of current stream", []string{"vhost", "app", "stream"}, prometheus.Labels{})
	promDescStatsStreamBytesInTotal  = prometheus.NewDesc("oven_stream_bytes_in_total", "total bytes in stream", []string{"vhost", "app", "stream"}, prometheus.Labels{})
	promDescStatsStreamBytesOutTotal = prometheus.NewDesc("oven_stream_bytes_out_total", "total bytes out stream", []string{"vhost", "app", "stream"}, prometheus.Labels{})
	promDescStatsStream              = []*prometheus.Desc{
		promDescStatsStreamClientTotal,
		promDescStatsStreamClientMax,
		promDescStatsStreamBytesInTotal,
		promDescStatsStreamBytesOutTotal,
	}
)

func ResponseStatsToMetrics(resp *api.ResponseStats, descs []*prometheus.Desc, labels ...string) []prometheus.Metric {
	if resp == nil || resp.Data == nil {
		return nil
	}
	list := []prometheus.Metric{}
	if m, err := prometheus.NewConstMetric(descs[0], prometheus.GaugeValue, float64(resp.Data.TotalConnections), labels...); err == nil {
		list = append(list, m)
	}
	if m, err := prometheus.NewConstMetric(descs[1], prometheus.GaugeValue, float64(resp.Data.MaxTotalConnections), labels...); err == nil {
		list = append(list, m)
	}
	if m, err := prometheus.NewConstMetric(descs[2], prometheus.GaugeValue, float64(resp.Data.TotalBytesIn), labels...); err == nil {
		list = append(list, m)
	}
	if m, err := prometheus.NewConstMetric(descs[3], prometheus.GaugeValue, float64(resp.Data.TotalBytesOut), labels...); err == nil {
		list = append(list, m)
	}
	return list
}

func (c *configData) Describe(d chan<- *prometheus.Desc) {
	for _, desc := range promDescStatsVHost {
		d <- desc
	}
	for _, desc := range promDescStatsApp {
		d <- desc
	}
	for _, desc := range promDescStatsStream {
		d <- desc
	}
}

func (c *configData) Collect(metrics chan<- prometheus.Metric) {
	respList, err := c.API.RequestListVHosts()
	if err != nil {
		log.Panicf("unable to fetch vhosts: %s", err)
	}
	for _, vhost := range respList.Data {
		logVhost := log.WithField("vhost", vhost)
		if resp, err := c.API.RequestStatsVHost(vhost); err == nil {
			for _, m := range ResponseStatsToMetrics(resp, promDescStatsVHost, vhost) {
				metrics <- m
			}
		}
		respList, err = c.API.RequestListApps(vhost)
		if err != nil {
			logVhost.Errorf("unable to fetch apps: %s", err)
			continue
		}
		for _, app := range respList.Data {
			logApp := logVhost.WithField("app", app)
			if resp, err := c.API.RequestStatsApp(vhost, app); err == nil {
				for _, m := range ResponseStatsToMetrics(resp, promDescStatsApp, vhost, app) {
					metrics <- m
				}
			}
			respList, err = c.API.RequestListStreams(vhost, app)
			if err != nil {
				logApp.Errorf("unable to fetch stream: %s", err)
				continue
			}
			for _, stream := range respList.Data {
				if resp, err := c.API.RequestStatsStream(vhost, app, stream); err == nil {
					for _, m := range ResponseStatsToMetrics(resp, promDescStatsStream, vhost, app, stream) {
						metrics <- m
					}
				}
			}
		}
	}
}
