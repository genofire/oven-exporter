package main

import (
	"github.com/bdlm/log"
)

func (config *configData) fetch() {
	respList, err := config.RequestListVHosts()
	if err != nil {
		log.Panicf("unable to fetch vhosts: %s", err)
	}
	for _, vhost := range respList.Data {
		logVhost := log.WithField("vhost", vhost)
		resp, err := config.RequestStatsVHost(vhost)
		if err != nil {
			logVhost.Errorf("error on request: %s", err)
		} else {
			resp.Log(logVhost)
		}
		respList, err = config.RequestListApps(vhost)
		if err != nil {
			logVhost.Errorf("unable to fetch apps: %s", err)
			continue
		}
		for _, app := range respList.Data {
			logApp := logVhost.WithField("app", app)
			resp, err = config.RequestStatsApp(vhost, app)
			if err != nil {
				logApp.Errorf("error on request: %s", err)
			} else {
				resp.Log(logApp)
			}
			respList, err = config.RequestListStreams(vhost, app)
			if err != nil {
				logApp.Errorf("unable to fetch stream: %s", err)
				continue
			}
			for _, stream := range respList.Data {
				logStream := logApp.WithField("stream", stream)
				req, err := config.RequestStatsStream(vhost, app, stream)
				if err != nil {
					logStream.Errorf("error on request: %s", err)
					continue
				}
				req.Log(logStream)
			}
		}
	}
}
