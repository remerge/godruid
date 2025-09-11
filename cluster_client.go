package godruid

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/remerge/rand"
)

type ClusterClient struct {
	clients []*Client
}

func NewDruidClusterClient(url string) (*ClusterClient, error) {
	return NewDruidClusterWithHTTP(url, http.DefaultClient)
}

func NewDruidClusterWithHTTP(url string, httpClient *http.Client) (*ClusterClient, error) {
	urls := strings.Split(url, ",")
	if len(urls) == 0 {
		return nil, errors.New("At least one druid url must be provided")
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	clients := make([]*Client, len(urls))
	for i := 0; i < len(urls); i++ {
		clients[i] = &Client{Url: urls[i], HttpClient: httpClient}
	}

	return &ClusterClient{clients}, nil
}

func (c *ClusterClient) sample() *Client {
	n := rand.Int() % len(c.clients)
	return c.clients[n]
}

func (c *ClusterClient) Query(query Query) (err error) {
	return c.sample().Query(query)
}

func (c *ClusterClient) QueryWithContext(ctx context.Context, query Query) (err error) {
	return c.sample().QueryWithContext(ctx, query)
}
