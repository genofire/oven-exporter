package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *configData) Request(url string, value interface{}) error {
	netClient := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(http.MethodGet, c.Host+url, nil)
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
