package openai

import (
	"errors"
	"net/http"
	"net/url"
	"time"
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

func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}
