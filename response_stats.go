package main

import (
	"fmt"

	"github.com/bdlm/log"
)

const (
	URLRequestStatsVHost  = "/v1/stats/current/vhosts/%s"
	URLRequestStatsApp    = "/v1/stats/current/vhosts/%s/apps/%s"
	URLRequestStatsStream = "/v1/stats/current/vhosts/%s/apps/%s/streams/%s"
)

type ResponseStats struct {
	Message    string             `json:"message"`
	StatusCode int                `json:"statusCode"`
	Data       *ResponseStatsData `json:"response,omitempty"`
}

type ResponseStatsData struct {
	// - timestamp - time.TIme has problem with nanosecond in JSON
	CreatedTime            string `json:"createdTime" example:"2021-07-19T23:13:12.162+0200"`
	LastRecvTime           string `json:"lastRecvTime" example:"2021-07-19T23:23:27.274+0200"`
	LastSentTime           string `json:"lastSentTime" example:"2021-07-19T23:23:27.077+0200"`
	LastUpdatedTime        string `json:"lastUpdatedTime" example:"2021-07-19T23:23:27.274+0200"`
	MaxTotalConnectionTime string `json:"maxTotalConnectionTime" example:"2021-07-19T23:16:37.851+0200"`
	// - coonnections
	TotalConnections    int `json:"totalConnections" example:"1"`
	MaxTotalConnections int `json:"maxTotalConnections" example:"2"`
	// - traffic
	TotalBytesIn  uint64 `json:"totalBytesIn" example:"120197570"`
	TotalBytesOut uint64 `json:"totalBytesOut" example:"117022184"`
}

func (resp *ResponseStats) Log(log *log.Entry) {
	logger := log
	if resp.Data != nil {
		logger.WithFields(map[string]interface{}{
			"max_clients": resp.Data.MaxTotalConnections,
			"clients":     resp.Data.TotalConnections,
		})
	}
	logger.Info(resp.Message)
}
func (c *configData) RequestStatsVHost(vhost string) (*ResponseStats, error) {
	req := ResponseStats{}
	url := fmt.Sprintf(URLRequestStatsVHost, vhost)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (c *configData) RequestStatsApp(vhost, app string) (*ResponseStats, error) {
	req := ResponseStats{}
	url := fmt.Sprintf(URLRequestStatsApp, vhost, app)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (c *configData) RequestStatsStream(vhost, app, stream string) (*ResponseStats, error) {
	req := ResponseStats{}
	url := fmt.Sprintf(URLRequestStatsStream, vhost, app, stream)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
