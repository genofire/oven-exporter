package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request to API and unmarshal result
func (c *Client) Request(method, url string, body, value interface{}) error {
	netClient := &http.Client{
		Timeout: time.Second * 20,
	}
	var jsonBody io.Reader
	if body != nil {
		if strBody, ok := body.(string); ok {
			jsonBody = strings.NewReader(strBody)
		} else {
			jsonBodyArray, err := json.Marshal(body)
			if err != nil {
				return err
			}
			jsonBody = bytes.NewBuffer(jsonBodyArray)
		}
	}
	req, err := http.NewRequest(method, c.URL+url, jsonBody)
	if err != nil {
		return err
	}
	req.Header = map[string][]string{
		"authorization": {fmt.Sprintf("Basic %s", c.Token)},
	}
	resp, err := netClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&value)
	resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}
