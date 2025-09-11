package godruid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	HttpClient *http.Client
}

func (c *Client) Query(query Query) error {
	return c.QueryWithContext(context.Background(), query)
}

func (c *Client) QueryWithContext(ctx context.Context, query Query) (err error) {
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
	result, err := c.QueryRawWithContext(ctx, reqJson)
	if err != nil {
		return
	}

	return query.OnResponse(result)
}

func (c *Client) QueryRaw(req []byte) ([]byte, error) {
	return c.QueryRawWithContext(context.Background(), req)
}

func (c *Client) QueryRawWithContext(ctx context.Context, req []byte) ([]byte, error) {
	if c.EndPoint == "" {
		c.EndPoint = DefaultEndPoint
	}
	endPoint := c.EndPoint
	if c.Debug {
		endPoint += "?pretty"
		c.LastRequest = string(req)
	}

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Url+endPoint, bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		c.LastResponse = string(payload)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", response.Status, string(payload))
	}

	return payload, nil
}
