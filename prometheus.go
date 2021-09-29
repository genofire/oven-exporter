package main

import (
	"dev.sum7.eu/genofire/oven-exporter/api"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
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

	promDescPushUp             = prometheus.NewDesc("oven_push_up", "state of push", []string{"vhost", "app", "stream", "id", "state"}, prometheus.Labels{})
	promDescPushSequence       = prometheus.NewDesc("oven_push_sequence", "sequence of started pushes", []string{"vhost", "app", "stream", "id"}, prometheus.Labels{})
	promDescPushSentBytes      = prometheus.NewDesc("oven_push_send_byte", "bytes send on push", []string{"vhost", "app", "stream", "id"}, prometheus.Labels{})
	promDescPushTotalSentBytes = prometheus.NewDesc("oven_push_total_send_bytes", "total bytes send on push", []string{"vhost", "app", "stream", "id"}, prometheus.Labels{})
	promDescPush               = []*prometheus.Desc{
		promDescPushUp,
		promDescPushSequence,
		promDescPushSentBytes,
		promDescPushTotalSentBytes,
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
	if m, err := prometheus.NewConstMetric(descs[1], prometheus.CounterValue, float64(resp.Data.MaxTotalConnections), labels...); err == nil {
		list = append(list, m)
	}
	if m, err := prometheus.NewConstMetric(descs[2], prometheus.CounterValue, float64(resp.Data.TotalBytesIn), labels...); err == nil {
		list = append(list, m)
	}
	if m, err := prometheus.NewConstMetric(descs[3], prometheus.CounterValue, float64(resp.Data.TotalBytesOut), labels...); err == nil {
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
	for _, desc := range promDescPush {
		d <- desc
	}
}

func (c *configData) Collect(metrics chan<- prometheus.Metric) {
	respList, err := c.API.RequestListVHosts()
	if err != nil {
		c.log.Panic("unable to fetch vhosts", zap.Error(err))
	}
	for _, vhost := range respList.Data {
		logVhost := c.log.With(zap.String("vhost", vhost))
		if resp, err := c.API.RequestStatsVHost(vhost); err == nil {
			for _, m := range ResponseStatsToMetrics(resp, promDescStatsVHost, vhost) {
				metrics <- m
			}
		}
		respList, err = c.API.RequestListApps(vhost)
		if err != nil {
			logVhost.Error("unable to fetch apps", zap.Error(err))
			continue
		}
		for _, app := range respList.Data {
			logApp := logVhost.With(zap.String("app", app))
			if resp, err := c.API.RequestStatsApp(vhost, app); err == nil {
				for _, m := range ResponseStatsToMetrics(resp, promDescStatsApp, vhost, app) {
					metrics <- m
				}
			}
			if resp, err := c.API.RequestPushStatus(vhost, app); err != nil {
				logApp.Error("unable to fetch pushes", zap.Error(err))
			} else {
				for _, data := range resp.Data {
					if m, err := prometheus.NewConstMetric(promDescPushUp, prometheus.GaugeValue, 1, vhost, app, data.Stream.Name, data.ID, data.State); err == nil {
						metrics <- m
					}
					labels := []string{vhost, app, data.Stream.Name, data.ID}
					if m, err := prometheus.NewConstMetric(promDescPushSequence, prometheus.CounterValue, float64(data.Sequence), labels...); err == nil {
						metrics <- m
					}
					if m, err := prometheus.NewConstMetric(promDescPushSentBytes, prometheus.GaugeValue, float64(data.SentBytes), labels...); err == nil {
						metrics <- m
					}
					if m, err := prometheus.NewConstMetric(promDescPushTotalSentBytes, prometheus.CounterValue, float64(data.TotalSentBytes), labels...); err == nil {
						metrics <- m
					}
				}
			}
			respList, err = c.API.RequestListStreams(vhost, app)
			if err != nil {
				logApp.Error("unable to fetch stream", zap.Error(err))
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
