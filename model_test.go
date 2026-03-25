package openai

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Models(t *testing.T) {
	client := testNewClient(t)

	models, err := client.Models(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, models)

	fmt.Println(models)
}
