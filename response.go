package openai

// ToolCall represents a tool call in the completion.
type ToolCall struct {
	// Type of the tool call, e.g., "function".
	Type string `json:"type"`

	// Unique identifier for the tool call.
	ID string `json:"id"`

	// Index of the tool call.
	Index int `json:"index"`

	// The function details for the call.
	Function *ToolCallFunction `json:"function"`
}

// ToolCallFunction represents a function call in the tool.
type ToolCallFunction struct {
	// Name of the function (required).
	Name string `json:"name"`

	// JSON string of arguments passed to the function (required).
	Arguments string `json:"arguments"`
}
