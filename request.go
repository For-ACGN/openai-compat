package openai

// ChatCompletionRequest defines the structure for a chat completion request.
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`                       // The ID of the model to use (required).
	Messages         []ChatCompletionMessage `json:"messages"`                    // A list of messages comprising the conversation (required).
	Audio            *ChatCompletionAudio    `json:"audio"`                       // audio output configuration.
	MaxTokens        int                     `json:"max_tokens,omitempty"`        // The maximum number of tokens to generate in the chat completion (optional).
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"` // Penalty for new tokens based on their frequency in the text so far (optional).
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`  // Penalty for new tokens based on their presence in the text so far (optional).
	Temperature      float32                 `json:"temperature,omitempty"`       // The sampling temperature, between 0 and 2 (optional).
	TopP             float32                 `json:"top_p,omitempty"`             // The nucleus sampling parameter, between 0 and 1 (optional).
	ResponseFormat   *ResponseFormat         `json:"response_format,omitempty"`   // The desired response format (optional).
	Stop             []string                `json:"stop,omitempty"`              // A list of sequences where the model should stop generating further tokens (optional).
	Tools            []Tool                  `json:"tools,omitempty"`             // A list of tools the model may use (optional).
	ToolChoice       any                     `json:"tool_choice,omitempty"`       // Controls which (if any) tool is called by the model (optional).
	LogProbs         bool                    `json:"logprobs,omitempty"`          // Whether to return log probabilities of the most likely tokens (optional).
	TopLogProbs      int                     `json:"top_logprobs,omitempty"`      // The number of top most likely tokens to return log probabilities for (optional).
}

func (ccr *ChatCompletionRequest) encode() []byte {
	return nil
}

// ChatCompletionMessage represents a single message in a chat completion conversation.
type ChatCompletionMessage struct {
	Role             string     `json:"role"`                        // The role of the message sender, e.g., "user", "assistant", "system".
	Content          string     `json:"content"`                     // The content of the message.
	ReasoningContent string     `json:"reasoning_content,omitempty"` // The reasoning content of the message (optional) when using the reasoner model with Chat Prefix Completion. When using this feature, the Prefix parameter must be set to true.
	ToolCallID       string     `json:"tool_call_id,omitempty"`      // Tool call that this message is responding to.
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`        // Optional tool calls.
}

// ChatCompletionAudio represents audio configuration in a chat completion conversation.
type ChatCompletionAudio struct {
	Format string `json:"format"` // set output audio format like wav, mp3.
	Voice  string `json:"voice"`  // set the voice of the model.
}

// ResponseFormat defines the structure for the response format.
type ResponseFormat struct {
	Type string `json:"type"` // The desired response format, either "text" or "json_object".
}

// Tool defines the structure for a tool.
type Tool struct {
	Type     string   `json:"type"`     // The type of the tool, e.g., "function" (required).
	Function Function `json:"function"` // The function details (required).
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
