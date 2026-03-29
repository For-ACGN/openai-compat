package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// CreateChatCompletion sends a chat completion request and returns the generated response.
func (c *Client) CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if req.Stream {
		return nil, errors.New("chat completion request is stream type")
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}
	hr, err := c.newPOSTRequest("", buf, false)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req.ctx, hr)
	if err != nil {
		return nil, err
	}
	defer c.closeBody(resp)
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, handleError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cr ChatCompletionResponse
	err = json.Unmarshal(body, &cr)
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

// CreateChatCompletionStream sends a chat completion request and returns the stream.
func (c *Client) CreateChatCompletionStream(req *ChatCompletionRequest) (*ChatCompletionStream, error) {
	if !req.Stream {
		return nil, errors.New("chat completion request is not stream type")
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, err
	}
	hr, err := c.newPOSTRequest("", buf, false)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req.ctx, hr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		c.closeBody(resp)
		return nil, handleError(resp)
	}

	stream := ChatCompletionStream{
		resp: resp,
	}
	return &stream, nil
}
