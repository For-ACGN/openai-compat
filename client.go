package openai

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// RoleDeveloper is the role of a developer message.
	RoleDeveloper = "developer"

	// RoleSystem is the role of a system message.
	RoleSystem = "system"

	// RoleUser is the role of a user message.
	RoleUser = "user"

	// RoleAssistant is the role of an assistant message.
	RoleAssistant = "assistant"

	// RoleTool is the role of a tool message.
	RoleTool = "tool"
)

const (
	defaultPath    = "chat/completions"
	defaultTimeout = 3 * time.Minute
)

// Client is a OpenAI compatibility client.
type Client struct {
	base string
	path string
	key  string

	client *http.Client
}

// Options contains the options for Client.
type Options struct {
	Path     string        `toml:"path"      json:"path"`
	Timeout  time.Duration `toml:"timeout"   json:"timeout"`
	ProxyURL string        `toml:"proxy_url" json:"proxy_url"`

	// if ProxyURL is not empty, it will cover field Proxy.
	Transport *http.Transport `toml:"-" json:"-"`
}

// NewClient is used to create a client.
func NewClient(baseURL, apiKey string, opts *Options) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("base url can not be empty")
	}
	if apiKey == "" {
		return nil, errors.New("api key can not be empty")
	}
	_, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if opts == nil {
		opts = new(Options)
	}
	path := opts.Path
	if path == "" {
		path = defaultPath
	}
	_, err = url.Parse(path)
	if err != nil {
		return nil, err
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	transport := opts.Transport
	if transport == nil {
		transport = new(http.Transport)
	}
	if opts.ProxyURL != "" {
		proxy, err := url.Parse(opts.ProxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = func(*http.Request) (*url.URL, error) {
			return proxy, nil
		}
	}
	client := Client{
		base: baseURL,
		path: path,
		key:  apiKey,
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
	}
	return &client, nil
}

func (c *Client) newPOSTRequest(path string, body io.Reader, stream bool) (*http.Request, error) {
	return c.newRequest(http.MethodPost, path, body, stream)
}

func (c *Client) newGETRequest(path string) (*http.Request, error) {
	return c.newRequest(http.MethodGet, path, nil, false)
}

func (c *Client) newRequest(method, path string, body io.Reader, stream bool) (*http.Request, error) {
	if path == "" {
		path = c.path
	}
	URL, err := url.JoinPath(c.base, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	if stream {
		req.Header.Set("cache-control", "no-cache")
	}
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.client.Do(req) // #nosec
}

func (c *Client) closeBody(resp *http.Response) {
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// Close is used to close connection in under http client.
func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}
