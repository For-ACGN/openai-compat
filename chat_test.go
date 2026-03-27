package openai

import (
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

const (
	testImageURL = "https://example-files.cnbj1.mi-fds.com/example-files/image/image_example.png"
	testAudioURL = "https://example-files.cnbj1.mi-fds.com/example-files/audio/audio_example.wav"
	testVideoURL = "https://example-files.cnbj1.mi-fds.com/example-files/video/video_example.mp4"
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

		testPrintResponse(resp)
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

		testPrintResponse(resp)
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

		testPrintResponse(resp)
	})

	t.Run("multi text content", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role: RoleUser,
				Content: []*Content{
					{Text: "Hello"},
					{Text: "LLM!"},
				},
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

		testPrintResponse(resp)
	})

	t.Run("with image content", func(t *testing.T) {
		t.Run("url", func(t *testing.T) {
			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "What is in this picture?"},
						{ImageURL: testImageURL},
					},
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
			require.Contains(t, response, "tree")
			require.Contains(t, response, "light")

			testPrintResponse(resp)
		})

		t.Run("data", func(t *testing.T) {
			image, err := os.ReadFile("testdata/image.png")
			require.NoError(t, err)

			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "What is in this picture?"},
						{ImageData: image},
					},
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
			require.Contains(t, response, "tree")
			require.Contains(t, response, "light")

			testPrintResponse(resp)
		})
	})

	t.Run("with audio content", func(t *testing.T) {
		t.Run("url", func(t *testing.T) {
			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "please describe the content of the audio"},
						{AudioURL: testAudioURL},
					},
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
			require.Contains(t, response, "day")
			require.Contains(t, response, "morning")

			testPrintResponse(resp)
		})

		t.Run("data", func(t *testing.T) {
			audio, err := os.ReadFile("testdata/audio.wav")
			require.NoError(t, err)

			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "please describe the content of the audio"},
						{AudioData: audio},
					},
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
			require.Contains(t, response, "day")
			require.Contains(t, response, "morning")

			testPrintResponse(resp)
		})
	})

	t.Run("with video content", func(t *testing.T) {
		t.Run("url", func(t *testing.T) {
			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "please describe the content of the video"},
						{VideoURL: testVideoURL, VideoFPS: 5, VideoResLevel: "default"},
					},
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
			require.Contains(t, response, "kitten")
			require.Contains(t, response, "blue")

			testPrintResponse(resp)
		})

		t.Run("data", func(t *testing.T) {
			video, err := os.ReadFile("testdata/video.mp4")
			require.NoError(t, err)

			req := NewChatCompletionRequest(false)
			req.Model = MiMoV2Omni
			req.Messages = []*ChatCompletionMessage{
				{
					Role:    RoleSystem,
					Content: "I'm writing a test, so please add prefix <test> in response",
				},
				{
					Role: RoleUser,
					Content: []*Content{
						{Text: "please describe the content of the video"},
						{VideoData: video, VideoFPS: 5, VideoResLevel: "default"},
					},
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
			require.Contains(t, response, "kitten")
			require.Contains(t, response, "blue")

			testPrintResponse(resp)
		})
	})

	err := client.Close()
	require.NoError(t, err)
}

func testPrintResponse(resp *ChatCompletionResponse) {
	fmt.Println(resp.ID)
	fmt.Println(resp.Model)
	fmt.Println(resp.Created)

	spew.Dump(resp.Choices)
	spew.Dump(resp.Usage)
}
