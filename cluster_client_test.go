package godruid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClusterClientCreation(t *testing.T) {
	url := "http://localhost:8081,http://localhosy:8082"

	cluster, err := NewDruidClusterClient(url)
	require.NoError(t, err)

	require.Len(t, cluster.clients, 2)
}
