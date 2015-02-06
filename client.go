package godruid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DefaultEndPoint = "/druid/v2"
)

type Client struct {
	Url      string
	EndPoint string

	Debug        bool
	LastRequest  string
	LastResponse string
}

func (c *Client) Query(query Query) (err error) {
	query.Setup()
	var reqJson []byte
	if c.Debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}
	if err != nil {
		return
	}
	result, err := c.QueryRaw(reqJson)
	if err != nil {
		return
	}

	return query.OnResponse(result)
}

func (c *Client) QueryRaw(req []byte) (result []byte, err error) {
	if c.EndPoint == "" {
		c.EndPoint = DefaultEndPoint
	}
	endPoint := c.EndPoint
	if c.Debug {
		endPoint += "?pretty"
		c.LastRequest = string(req)
	}
	if err != nil {
		return
	}

	resp, err := http.Post(c.Url+endPoint, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return
	}
	defer func() {
		resp.Body.Close()
	}()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if c.Debug {
		c.LastResponse = string(result)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(result))
	}

	return
}
