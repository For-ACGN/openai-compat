package openai

import (
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
}
