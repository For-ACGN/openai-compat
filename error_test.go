package openai

import (
	"fmt"
	"testing"
	
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	client := testNewClient(t)
	
	req, err := client.newGETRequest("models_invalid")
	require.NoError(t, err)
	
	resp, err := client.do(req)
	require.NoError(t, err)
	
	apiErr := handleError(resp)
	ae, ok := apiErr.(*APIError)
	require.True(t, ok)
	fmt.Println(ae)
	
	err = client.Close()
	require.NoError(t, err)
}
