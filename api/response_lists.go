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

// RequestListVHosts to get list of vhosts
func (c *Client) RequestListVHosts() (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListVHost)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// RequestListApps to get list of apps on given vhost
func (c *Client) RequestListApps(vhost string) (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListApp, vhost)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// RequestDefaultListApps to get list of apps on default vhost
func (c *Client) RequestDefaultListApps() (*ResponseList, error) {
	return c.RequestListApps(c.DefaultVHost)
}

// RequestDefaultListStreams to get list of streams on given vhost and app
func (c *Client) RequestListStreams(vhost, app string) (*ResponseList, error) {
	req := ResponseList{}
	url := fmt.Sprintf(URLRequestListStream, vhost, app)
	if err := c.Request(url, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// RequestDefaultListStreams to get list of streams on default vhost and app
func (c *Client) RequestDefaultListStreams() (*ResponseList, error) {
	return c.RequestListStreams(c.DefaultVHost, c.DefaultApp)
}
