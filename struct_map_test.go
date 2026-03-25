package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapStruct(t *testing.T) {
	type testA struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	ta := testA{
		A: "test",
		B: 123,
	}
	m, err := mapStruct(ta)
	require.NoError(t, err)
	fmt.Println(m)

	require.Equal(t, m["a"], "test")
	require.Equal(t, m["b"], float64(123))

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(m)
	require.NoError(t, err)
	fmt.Println(buf.String())
	require.Contains(t, buf.String(), "\"a\"")
}
