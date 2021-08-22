package api

import (
	"fmt"
	"net/http"
)

// API URLS for Push
const (
	URLRequestPushStatus = "/v1/vhosts/%s/apps/%s:pushes"
	URLRequestPushStart  = "/v1/vhosts/%s/apps/%s:startPush"
	URLRequestPushStop   = "/v1/vhosts/%s/apps/%s:stopPush"
)

// ResponsePushStatus JSON Message with Push status data
type ResponsePushStatus struct {
	Message    string              `json:"message"`
	StatusCode int                 `json:"statusCode"`
	Data       []*ResponsePushData `json:"response,omitempty"`
}

// ResponsePushData one push configuration
type ResponsePushData struct {
	VHost     string                  `json:"vhost" example:"default"`
	App       string                  `json:"app" example:"live"`
	ID        string                  `json:"id" example:"youtube"`
	Stream    *ResponsePushDataStream `json:"stream"`
	State     string                  `json:"state" example:"ready"`
	Protocol  string                  `json:"protocol" example:"rtmp"`
	URL       string                  `json:"url" example:"rtmp://a.rtmp.youtube.com/live2"`
	StreamKey string                  `json:"streamKey" example:"SUPERSECRET"`
	// - timestamp - time.TIme has problem with nanosecond in JSON
	CreatedTime  string `json:"createdTime" example:"2021-07-19T23:13:12.162+0200"`
	FinishedTime string `json:"finishedTime" example:"2021-07-19T23:23:27.274+0200"`
	StartTime    string `json:"startTime" example:"2021-07-19T23:23:27.077+0200"`
	// - coonnections
	Sequence int `json:"sequence" example:"1"`
	// - traffic
	SentBytes      uint64 `json:"sentBytes" example:"0"`
	SentTime       uint64 `json:"sentTime" example:"0"`
	TotalSentBytes uint64 `json:"totalsentBytes" example:"356233652"`
	TotalSentTime  uint64 `json:"totalsentTime" example:"933808"`
}

// ResponsePushStart one push configuration
type ResponsePushStart struct {
	ID        string                  `json:"id" example:"youtube"`
	Stream    *ResponsePushDataStream `json:"stream"`
	Protocol  string                  `json:"protocol" example:"rtmp"`
	URL       string                  `json:"url" example:"rtmp://a.rtmp.youtube.com/live2"`
	StreamKey string                  `json:"streamKey" example:"SUPERSECRET"`
}

// ResponsePushDataStream of data of stream
type ResponsePushDataStream struct {
	Name   string `json:"name" example:"..."`
	Tracks []int  `json:"tracks" example:"[]"`
}

// RequestPushStatus to get list of pushes and his configuration
func (c *Client) RequestPushStatus(vhost, app string) (*ResponsePushStatus, error) {
	req := ResponsePushStatus{}
	url := fmt.Sprintf(URLRequestPushStatus, vhost, app)
	if err := c.Request(http.MethodGet, url, nil, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// RequestPushStatusDefault to get list of pushes and his configuration for default vhost and app
func (c *Client) RequestPushStatusDefault() (*ResponsePushStatus, error) {
	return c.RequestPushStatus(c.DefaultVHost, c.DefaultApp)
}

// StartPush to delete an push
func (c *Client) StartPush(vhost, app string, data *ResponsePushStart) (*ResponsePushStatus, error) {
	req := ResponsePushStatus{}
	url := fmt.Sprintf(URLRequestPushStart, vhost, app)
	if err := c.Request(http.MethodPost, url, &data, &req); err != nil {
		return &req, err
	}
	return &req, nil
}

// StartPushDefault to delete an push and on default vhost and app
func (c *Client) StartPushDefault(data *ResponsePushStart) (*ResponsePushStatus, error) {
	return c.StartPush(c.DefaultVHost, c.DefaultApp, data)
}

// DeletePush to delete an push
func (c *Client) DeletePush(vhost, app, id string) error {
	type idJSON struct {
		ID string `json:"id"`
	}
	url := fmt.Sprintf(URLRequestPushStop, vhost, app)
	data := idJSON{ID: id}
	if err := c.Request(http.MethodPost, url, &data, nil); err != nil {
		return err
	}
	return nil
}

// DeletePushDefault to delete an push and on default vhost and app
func (c *Client) DeletePushDefault(id string) error {
	return c.DeletePush(c.DefaultVHost, c.DefaultApp, id)
}
