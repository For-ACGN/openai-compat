package openai

import (
	"encoding/json"
	"maps"
)

// ChatCompletionRequest defines the structure for a chat completion request.
type ChatCompletionRequest struct {
	// The ID of the model to use (required).
	Model string `json:"model"`

	// A list of messages comprising the conversation (required).
	Messages []*ChatCompletionMessage `json:"messages"`

	// audio output configuration.
	Audio *ChatCompletionAudio `json:"audio,omitempty"`

	// A list of tools the model may use (optional).
	Tools []*Tool `json:"tools,omitempty"`

	// Controls which (if any) tool is called by the model (optional).
	ToolChoice any `json:"tool_choice,omitempty"`

	// controlling the transition between thinking and non-thinking modes
	Thinking bool `json:"thinking,omitempty"`

	// A list of sequences where the model should stop generating further tokens (optional).
	Stop []string `json:"stop,omitempty"`

	// The maximum number of tokens to generate in the chat completion (optional).
	MaxTokens int `json:"max_tokens,omitempty"`

	// Penalty for new tokens based on their frequency in the text so far (optional).
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`

	// Penalty for new tokens based on their presence in the text so far (optional).
	PresencePenalty float32 `json:"presence_penalty,omitempty"`

	// The sampling temperature, between 0 and 2 (optional).
	Temperature float32 `json:"temperature,omitempty"`

	// The nucleus sampling parameter, between 0 and 1 (optional).
	TopP float32 `json:"top_p,omitempty"`

	// The desired response format (optional).
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`

	// To set additional parameters for a specific model, note that
	// the hierarchy is top-level with ChatCompletionRequest.
	Extra map[string]any `json:"-"`

	// stream internal switch.
	stream bool
}

func (ccr *ChatCompletionRequest) encode() ([]byte, error) {
	mv, err := mapStruct(ccr)
	if err != nil {
		return nil, err
	}
	// process argument about thinking
	thinking := new(Thinking)
	if ccr.Thinking {
		thinking.Type = "enabled"
	} else {
		thinking.Type = "disabled"
	}
	mv["thinking"] = thinking
	// process different fields about max_tokens.
	if ccr.MaxTokens != 0 {
		switch ccr.Model {
		case MiMoV2Flash, MiMoV2Omni, MiMoV2Pro, MiMoV2TTS:
			delete(mv, "max_tokens")
			mv["max_completion_tokens"] = ccr.MaxTokens
		}
	}
	// set stream mode
	mv["stream"] = ccr.stream
	// append extra arguments
	maps.Copy(mv, ccr.Extra)
	return json.Marshal(mv)
}

// ChatCompletionMessage represents a single message in a chat completion conversation.
type ChatCompletionMessage struct {
	// The role of the message sender, e.g., "user", "assistant", "system".
	Role string `json:"role"`

	// The content of the message, the type can be string or []*Content.
	Content any `json:"content"`

	// Tool call that this message is responding to.
	ToolCallID string `json:"tool_call_id,omitempty"`

	// Optional tool calls.
	ToolCalls []*ToolCall `json:"tool_calls,omitempty"`

	// Optional names for participants. Provides information for
	// the model to distinguish participants with the same role.
	Name string `json:"name,omitempty"`
}

// Content is set to ChatCompletionMessage field.
type Content struct {
	// the type of content like "text", "image_url", "input_audio", "video_url".
	Type string `json:"type"`

	// text data
	Text string `json:"text,omitempty"`

	// image url or base64 encoded data.
	ImageURL *ImageURL `json:"image_url,omitempty"`

	// audio url or base64 encoded data.
	InputAudio *InputAudio `json:"input_audio,omitempty"`

	// video url or base64 encoded data.
	VideoURL *VideoURL `json:"video_url,omitempty"`

	// video frames per second.
	FPS int `json:"fps,omitempty"`

	// video Resolution level.
	MediaResolution string `json:"media_resolution,omitempty"`
}

// MarshalJSON implement interface json.Marshaler.
func (c *Content) MarshalJSON() ([]byte, error) {
	type alias Content
	tmp := alias(*c)
	mv, err := mapStruct(tmp)
	if err != nil {
		return nil, err
	}
	return json.Marshal(mv)
}

// ImageURL for content.
type ImageURL struct {
	URL string `json:"url"`
}

// InputAudio for content.
type InputAudio struct {
	Data string `json:"data"`
}

// VideoURL for content.
type VideoURL struct {
	URL string `json:"url"`
}

// ChatCompletionAudio represents audio configuration in a chat completion conversation.
type ChatCompletionAudio struct {
	Format string `json:"format"` // set output audio format like wav, mp3.
	Voice  string `json:"voice"`  // set the voice of the model.
}

// Tool defines the structure for a tool.
type Tool struct {
	// The type of the tool, e.g., "function" (required).
	Type string `json:"type"`

	// The function details (required).
	Function *Function `json:"function"`
}

// Function defines the structure of a function tool.
type Function struct {
	// The name of the function (required).
	Name string `json:"name"`

	// A description of the function (required).
	Description string `json:"description"`

	// The parameters of the function (optional).
	Parameters *FunctionParameters `json:"parameters,omitempty"`
}

// FunctionParameters defines the parameters for a function.
type FunctionParameters struct {
	// The type of the parameters, e.g., "object" (required).
	Type string `json:"type"`

	// The properties of the parameters (optional).
	Properties map[string]any `json:"properties,omitempty"`

	// A list of required parameter names (optional).
	Required []string `json:"required,omitempty"`
}

// ToolChoice defines the structure for a tool choice.
type ToolChoice struct {
	// The type of the tool, e.g., "function" (required).
	Type string `json:"type"`

	// The function details (optional, but required if type is "function").
	Function *ToolChoiceFunction `json:"function,omitempty"`
}

// ToolChoiceFunction defines the function details within ToolChoice.
type ToolChoiceFunction struct {
	// The name of the function to call (required).
	Name string `json:"name"`
}

// Thinking is used to control enable reasoning.
type Thinking struct {
	Type string `json:"type"`
}

// ResponseFormat defines the structure for the response format.
type ResponseFormat struct {
	// The desired response format, either "text" or "json_object".
	Type string `json:"type"`
}
