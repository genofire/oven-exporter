package api

import "encoding/base64"

// A Client for the API
type Client struct {
	Token string `toml:"token"`
	Host  string `toml:"host"`
}

// New Client from host and token
func New(host, token string) *Client {
	c := &Client{
		Host: host,
	}
	c.SetToken(token)
	return c
}

// SetToken by using base64encoding
func (c *Client) SetToken(token string) {
	c.Token = base64.StdEncoding.EncodeToString([]byte(token))
}
