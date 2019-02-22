package godruid

import (
	"errors"
	"strings"

	"github.com/remerge/rand"
)

type ClusterClient struct {
	clients []*Client
}

func NewDruidClusterFromUrl(url string) (*ClusterClient, error) {
	urls := strings.Split(url, ",")
	if len(urls) == 0 {
		return nil, errors.New("At least one druid url must be provided")
	}

	clients := make([]*Client, len(urls))
	for i := 0; i < len(urls); i++ {
		clients[i] = &Client{Url: urls[i]}
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
