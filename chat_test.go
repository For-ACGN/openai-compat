package openai

import (
	"fmt"
	"io"
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

var testToolGetLocation = &Function{
	model: MiMoV2Omni,

	Name:        "GetLocation",
	Description: "get user current location",
}

var testToolGetTemperature = &Function{
	model: MiMoV2Omni,

	Name:        "GetTemperature",
	Description: "get temperature by city name",
	Parameters: &FunctionParameters{
		Type: "object",
		Properties: map[string]*Property{
			"city": {Type: "string", Description: "input city name"},
		},
		Required: []string{"city"},
	},
}

var testToolGetRelativeHumidity = &Function{
	model: MiMoV2Omni,

	Name:        "GetRelativeHumidity",
	Description: "get relative humidity by city name",
	Parameters: &FunctionParameters{
		Type: "object",
		Properties: map[string]*Property{
			"city": {Type: "string", Description: "input city name"},
		},
		Required: []string{"city"},
	},
}

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

	t.Run("tool call", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role:    RoleUser,
				Content: "What is the current temperature and relative humidity?",
			},
		}
		req.Tools = []any{
			testToolGetLocation,
			testToolGetTemperature,
			testToolGetRelativeHumidity,
		}

		resp, err := client.CreateChatCompletion(req)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		message := resp.Choices[0].Message
		response := message.Content
		fmt.Println(response)
		require.Contains(t, response, "temperature")

		require.Equal(t, RoleAssistant, message.Role)
		toolCalls := message.ToolCalls
		require.NotEmpty(t, toolCalls)

		toolCall := toolCalls[0]
		require.Equal(t, "function", toolCall.Type)
		require.NotEmpty(t, toolCall.ID)

		fn := toolCalls[0].Function
		require.Equal(t, "GetLocation", fn.Name)
		require.Equal(t, "{}", fn.Arguments)

		question := &ChatCompletionMessage{
			Role:       RoleAssistant,
			Content:    message.Content,
			ToolCallID: toolCall.ID,
			ToolCalls:  message.ToolCalls,
		}
		req.Messages = append(req.Messages, question)
		req.Messages = append(req.Messages, &ChatCompletionMessage{
			Role:       RoleTool,
			Content:    "the current location is Shanghai",
			ToolCallID: toolCall.ID,
			ToolCalls:  toolCalls,
		})

		resp, err = client.CreateChatCompletion(req)
		require.NoError(t, err)

		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		message = resp.Choices[0].Message
		response = message.Content
		fmt.Println(response)
		require.Contains(t, response, "Shanghai")

		require.Len(t, message.ToolCalls, 2)
		toolCall1 := toolCalls[0]
		require.Equal(t, "function", toolCall1.Type)
		require.NotEmpty(t, toolCall1.ID)

		toolCall2 := toolCalls[0]
		require.Equal(t, "function", toolCall2.Type)
		require.NotEmpty(t, toolCall2.ID)

		question = &ChatCompletionMessage{
			Role:       RoleAssistant,
			Content:    message.Content,
			ToolCallID: toolCall1.ID,
			ToolCalls:  message.ToolCalls,
		}
		req.Messages = append(req.Messages, question)
		req.Messages = append(req.Messages, &ChatCompletionMessage{
			Role:       RoleTool,
			Content:    "the temperature is 25°C",
			ToolCallID: toolCall1.ID,
			ToolCalls:  toolCalls,
		})

		question = &ChatCompletionMessage{
			Role:       RoleAssistant,
			Content:    message.Content,
			ToolCallID: toolCall2.ID,
			ToolCalls:  message.ToolCalls,
		}
		req.Messages = append(req.Messages, question)
		req.Messages = append(req.Messages, &ChatCompletionMessage{
			Role:       RoleTool,
			Content:    "the relative humidity is 50%",
			ToolCallID: toolCall2.ID,
			ToolCalls:  toolCalls,
		})

		resp, err = client.CreateChatCompletion(req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.ID)
		require.Equal(t, MiMoV2Omni, resp.Model)
		require.NotEmpty(t, resp.Choices)
		require.NotZero(t, resp.Usage)
		require.NotZero(t, resp.Created)

		message = resp.Choices[0].Message
		response = message.Content
		fmt.Println(response)

		require.Contains(t, response, "25")
		require.Contains(t, response, "50%")
	})

	t.Run("web search", func(t *testing.T) {
		req := NewChatCompletionRequest(false)
		req.Model = MiMoV2Omni
		req.Messages = []*ChatCompletionMessage{
			{
				Role:    RoleSystem,
				Content: "I'm writing a test, so please add prefix <test> in response",
			},
			{
				Role:    RoleUser,
				Content: "What is the current temperature in Shanghai",
			},
		}
		req.Tools = []any{
			NewWebSearchTool(MiMoV2Omni),
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
		require.Contains(t, response, "°C")
	})

	err := client.Close()
	require.NoError(t, err)
}

func TestClient_CreateChatCompletionStream(t *testing.T) {
	client := testNewClient(t)

	t.Run("common", func(t *testing.T) {
		req := NewChatCompletionRequest(true)
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

		resp, err := client.CreateChatCompletionStream(req)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, resp.Close())
		}()

		var response string
		for {
			stream, err := resp.Receive()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			require.NotEmpty(t, stream.ID)
			require.Equal(t, MiMoV2Omni, stream.Model)
			require.NotZero(t, stream.Created)

			if len(stream.Choices) > 0 {
				delta := stream.Choices[0].Delta.Content
				fmt.Println(delta)
				response += delta
			}
		}
		require.Contains(t, response, "<test>")
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
