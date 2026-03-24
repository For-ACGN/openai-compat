package openai

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testBaseURL = "https://api.xiaomimimo.com/v1"
)

func testNewClient(t *testing.T) *Client {
	apiKey := os.Getenv("API_KEY")
	opts := Options{
		Timeout: 30 * time.Second,
	}
	client, err := NewClient(testBaseURL, apiKey, &opts)
	require.NoError(t, err)
	return client
}

func TestNewClient(t *testing.T) {
	client := testNewClient(t)

	err := client.Close()
	require.NoError(t, err)
}
