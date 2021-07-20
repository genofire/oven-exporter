package main

import (
	"dev.sum7.eu/genofire/oven-exporter/api"
	"github.com/bdlm/log"
)

func statsLog(resp *api.ResponseStats, log *log.Entry) {
	logger := log
	if resp.Data != nil {
		logger = logger.WithFields(map[string]interface{}{
			"max_clients": resp.Data.MaxTotalConnections,
			"clients":     resp.Data.TotalConnections,
		})
	}
	logger.Info(resp.Message)
}

func fetch(client *api.Client) {
	respList, err := client.RequestListVHosts()
	if err != nil {
		log.Panicf("unable to fetch vhosts: %s", err)
	}
	for _, vhost := range respList.Data {
		logVhost := log.WithField("vhost", vhost)
		resp, err := client.RequestStatsVHost(vhost)
		if err != nil {
			logVhost.Errorf("error on request: %s", err)
		} else {
			statsLog(resp, logVhost)
		}
		respList, err = client.RequestListApps(vhost)
		if err != nil {
			logVhost.Errorf("unable to fetch apps: %s", err)
			continue
		}
		for _, app := range respList.Data {
			logApp := logVhost.WithField("app", app)
			resp, err = client.RequestStatsApp(vhost, app)
			if err != nil {
				logApp.Errorf("error on request: %s", err)
			} else {
				statsLog(resp, logApp)
			}
			respList, err = client.RequestListStreams(vhost, app)
			if err != nil {
				logApp.Errorf("unable to fetch stream: %s", err)
				continue
			}
			for _, stream := range respList.Data {
				logStream := logApp.WithField("stream", stream)
				req, err := client.RequestStatsStream(vhost, app, stream)
				if err != nil {
					logStream.Errorf("error on request: %s", err)
					continue
				}
				statsLog(req, logStream)
			}
		}
	}
}
