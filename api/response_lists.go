package api

import (
	"fmt"
)

const (
	URLRequestListVHost  = "/v1/vhosts"
	URLRequestListApp    = "/v1/vhosts/%s/apps"
	URLRequestListStream = "/v1/vhosts/%s/apps/%s/streams"
)

type ResponseList struct {
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
	Data       []string `json:"response,omitempty"`
}

func (c *Client) RequestListVHosts() (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListVHost)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (c *Client) RequestListApps(vhost string) (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListApp, vhost)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (c *Client) RequestListStreams(vhost, app string) (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListStream, vhost, app)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
