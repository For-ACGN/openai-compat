package openai

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateChatCompletion(t *testing.T) {
	client := testNewClient(t)

	t.Run("common", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role:    RoleUser,
				Content: "Hello LLM!",
			},
		}

		resp, err := client.CreateChatCompletion(req)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		response := resp.Choices[0].Message.Content
		fmt.Println(response)
		require.Contains(t, response, "<test>")

		fmt.Println(resp.ID)
		fmt.Println(resp.Model)
		fmt.Println(resp.Created)

		spew.Dump(resp.Choices)
		spew.Dump(resp.Usage)
	})

	t.Run("thinking", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role:    RoleUser,
				Content: "Hello LLM!",
			},
		}
		req.Thinking = true

		resp, err := client.CreateChatCompletion(req)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		response := resp.Choices[0].Message.Content
		fmt.Println(response)
		require.Contains(t, response, "<test>")

		reason := resp.Choices[0].Message.ReasoningContent
		fmt.Println(reason)
		require.NotEmpty(t, reason)

		fmt.Println(resp.ID)
		fmt.Println(resp.Model)
		fmt.Println(resp.Created)

		spew.Dump(resp.Choices)
		spew.Dump(resp.Usage)
	})

	t.Run("max tokens", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role:    RoleUser,
				Content: "Hello LLM!",
			},
		}
		req.MaxTokens = 5

		resp, err := client.CreateChatCompletion(req)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		response := resp.Choices[0].Message.Content
		fmt.Println(response)
		require.Contains(t, response, "<test>")
		require.Less(t, len(response), 20)

		fmt.Println(resp.ID)
		fmt.Println(resp.Model)
		fmt.Println(resp.Created)

		spew.Dump(resp.Choices)
		spew.Dump(resp.Usage)
	})

	err := client.Close()
	require.NoError(t, err)
}
