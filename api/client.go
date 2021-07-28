package api

import "encoding/base64"

// A Client for the API
type Client struct {
	Token        string `toml:"token"`
	URL          string `toml:"url"`
	DefaultVHost string `toml:"default_vhost"`
	DefaultApp   string `toml:"default_app"`
}

// New Client from host and token
func New(url, token string) *Client {
	c := &Client{
		URL: url,
	}
	c.SetToken(token)
	return c
}

// SetToken by using base64encoding
func (c *Client) SetToken(token string) {
	c.Token = base64.StdEncoding.EncodeToString([]byte(token))
}
