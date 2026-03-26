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
	Role             string      `json:"role"`                        // The role of the message sender, e.g., "user", "assistant", "system".
	Content          string      `json:"content"`                     // The content of the message.
	ReasoningContent string      `json:"reasoning_content,omitempty"` // The reasoning content of the message (optional) when using the reasoner model with Chat Prefix Completion. When using this feature, the Prefix parameter must be set to true.
	ToolCallID       string      `json:"tool_call_id,omitempty"`      // Tool call that this message is responding to.
	ToolCalls        []*ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls.
}

// ChatCompletionAudio represents audio configuration in a chat completion conversation.
type ChatCompletionAudio struct {
	Format string `json:"format"` // set output audio format like wav, mp3.
	Voice  string `json:"voice"`  // set the voice of the model.
}

// Tool defines the structure for a tool.
type Tool struct {
	Type     string    `json:"type"`     // The type of the tool, e.g., "function" (required).
	Function *Function `json:"function"` // The function details (required).
}

// ToolChoice defines the structure for a tool choice.
type ToolChoice struct {
	Type     string             `json:"type"`               // The type of the tool, e.g., "function" (required).
	Function ToolChoiceFunction `json:"function,omitempty"` // The function details (optional, but required if type is "function").
}

// ToolChoiceFunction defines the function details within ToolChoice.
type ToolChoiceFunction struct {
	// The name of the function to call (required).
	Name string `json:"name"`
}

// Function defines the structure of a function tool.
type Function struct {
	Name        string              `json:"name"`                 // The name of the function (required).
	Description string              `json:"description"`          // A description of the function (required).
	Parameters  *FunctionParameters `json:"parameters,omitempty"` // The parameters of the function (optional).
}

// FunctionParameters defines the parameters for a function.
type FunctionParameters struct {
	Type       string         `json:"type"`                 // The type of the parameters, e.g., "object" (required).
	Properties map[string]any `json:"properties,omitempty"` // The properties of the parameters (optional).
	Required   []string       `json:"required,omitempty"`   // A list of required parameter names (optional).
}

// Thinking is used to control enable reasoning.
type Thinking struct {
	Type string `json:"type"`
}

// ResponseFormat defines the structure for the response format.
type ResponseFormat struct {
	Type string `json:"type"` // The desired response format, either "text" or "json_object".
}
