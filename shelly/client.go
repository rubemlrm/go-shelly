package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-cleanhttp"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// Client act's as the entry object for sdk
type Client struct {
	BaseURL     *url.URL
	client      *retryablehttp.Client
	common      service
	username    string
	password    string
	requireAuth bool
}

type service struct {
	client *Client
}

type Response struct {
	*http.Response
}

// NewClient creates a new http client instance in case the provided one is nil
func NewClient(hostname string) (*Client, error) {
	client, err := newClient(hostname)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewBasicAuthClient(hostname, username, password string, requireAuth bool) (*Client, error) {
	client, err := newClient(hostname)
	if err != nil {
		return nil, err
	}

	client.username = username
	client.password = password
	client.requireAuth = requireAuth
	return client, nil
}

func newClient(hostname string) (*Client, error) {
	baseURL, err := url.Parse(hostname)
	if err != nil {
		return nil, err
	}

	c := &Client{}
	c.BaseURL = baseURL
	// Configure the HTTP client.
	c.client = &retryablehttp.Client{
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 400 * time.Millisecond,
		RetryMax:     5,
		CheckRetry:   c.retryHTTPCheck,
	}
	return c, nil
}

func (c *Client) retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return false, err
	}
	if resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}

func (c *Client) NewRequest(method, endpoint string, parseResponse interface{}) (*retryablehttp.Request, error) {
	composedURL := fmt.Sprintf("%v%v", c.BaseURL, endpoint)
	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")

	if method == http.MethodPatch || method == http.MethodPost || method == http.MethodPut {
		reqHeaders.Set("Content-Type", "application/json")
	}

	request, err := retryablehttp.NewRequest(method, composedURL, nil)
	if c.username != "" {
		request.SetBasicAuth(c.username, c.password)
	}
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (c *Client) SetBasicAuth(request *retryablehttp.Request) error {
	if c.username == "" {
		return fmt.Errorf("username can't be empty")
	}
	if c.password == "" {
		return fmt.Errorf("password can't be empty")
	}
	request.SetBasicAuth(c.username, c.password)
	return nil
}

func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 500 {
		return nil, errors.New("Server error")
	}

	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)
	_ = &Response{Response: resp}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return nil, nil
}
