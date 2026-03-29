package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Official DeepSeek Models.
const (
	DeepSeekChat     = "deepseek-chat"
	DeepSeekReasoner = "deepseek-reasoner"
)

// Official Xiaomi MiMo Models.
const (
	MiMoV2Flash = "mimo-v2-flash"
	MiMoV2Omni  = "mimo-v2-omni"
	MiMoV2Pro   = "mimo-v2-pro"
	MiMoV2TTS   = "mimo-v2-tts"
)

// Model is a model that can be used with the API.
type Model struct {
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
	ID      string `json:"id"`
}

// APIModels contains the response from the API endpoint.
type APIModels struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// Models is used to get the available models.
func (c *Client) Models(ctx context.Context) ([]string, error) {
	req, err := c.newGETRequest("models")
	if err != nil {
		return nil, err
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer c.closeBody(resp)
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, handleAPIError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var models APIModels
	err = json.Unmarshal(body, &models)
	if err != nil {
		return nil, err
	}

	var list []string
	for i := 0; i < len(models.Data); i++ {
		list = append(list, models.Data[i].ID)
	}
	return list, nil
}
