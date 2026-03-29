package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ChatCompletionStream is a response stream.
type ChatCompletionStream struct {
	resp *http.Response
}

// Receive is used to receive delta data.
func (s *ChatCompletionStream) Receive() (*ChatCompletionStreamResponse, error) {
	reader := bufio.NewReader(s.resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, fmt.Errorf("failed to read stream: %s", err)
		}
		if line == "data: [DONE]" {
			return nil, io.EOF
		}
		if len(line) > 6 && line[:6] == "data: " {
			trimmed := []byte(line[6:])
			var resp ChatCompletionStreamResponse
			err = json.Unmarshal(trimmed, &resp)
			if err != nil {
				return nil, fmt.Errorf("unmarshal error: %s, raw data: %s", err, trimmed)
			}
			if resp.Usage == nil {
				resp.Usage = new(Usage)
			}
			return &resp, nil
		}
	}
}

// Close is used to close the response body.
func (s *ChatCompletionStream) Close() error {
	return s.resp.Body.Close()
}

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
