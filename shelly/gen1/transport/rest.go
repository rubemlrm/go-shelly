package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	contracts "github.com/rubemlrm/go-shelly/shelly/gen1/contracts"
)

type ClientProxy interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

var _ contracts.ShellyClient = (*Client)(nil)

type ClientOptions struct {
	Hostname string
	Username string
	Password string
}

// Client act's as the entry object for sdk
type Client struct {
	client       ClientProxy
	Password     string
	Username     string
	RequiresAuth bool
	BaseURL      *url.URL
}

// NewClient creates a new http client instance in case the provided one is nil
func NewRestClient(options ClientOptions) (*Client, error) {
	client, err := newClient(options.Hostname)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRestBasicAuthClient(opts ClientOptions) (*Client, error) {
	client, err := newClient(opts.Hostname)
	if err != nil {
		return nil, err
	}

	client.Username = opts.Username
	client.Password = opts.Password
	client.RequiresAuth = true
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
		CheckRetry:   c.RetryHTTPCheck,
	}
	return c, nil
}

func (c *Client) RetryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
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

func (c *Client) NewRequest(method, endpoint string, opts interface{}) (*retryablehttp.Request, error) {
	jsonMethodsList := []string{
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	}
	var body interface{}

	u, err := c.ParseUrl(method, endpoint, opts)
	if err != nil {
		return nil, err
	}

	reqHeaders := make(http.Header)
	reqHeaders.Set("Accept", "application/json")
	if slices.Contains(jsonMethodsList, method) {
		reqHeaders.Set("Content-Type", "application/json")
		if opts != nil {
			body, err = json.Marshal(opts)
			if err != nil {
				return nil, err
			}
		}
	}

	request, err := retryablehttp.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	c.SetAdditionalHeaders(request, reqHeaders)

	if c.RequiresAuth {
		err = c.SetBasicAuth(request)
		if err != nil {
			return nil, err
		}
	}

	return request, nil
}

func (c *Client) ParseUrl(method, endpoint string, opts interface{}) (string, error) {
	u := *c.BaseURL
	unescaped, err := url.PathUnescape(endpoint)
	if err != nil {
		return "", err
	}

	// Set the encoded path data
	u.RawPath = c.BaseURL.Path + endpoint
	u.Path = c.BaseURL.Path + unescaped

	if method == http.MethodGet && opts != nil {
		q, err := query.Values(opts)
		if err != nil {
			return "", err
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func (c *Client) SetAdditionalHeaders(request *retryablehttp.Request, headers http.Header) {
	for k, v := range headers {
		request.Header[k] = v
	}
}

func (c *Client) SetBasicAuth(request *retryablehttp.Request) error {
	if c.Username == "" {
		return fmt.Errorf("username can't be empty")
	}
	if c.Password == "" {
		return fmt.Errorf("password can't be empty")
	}
	request.SetBasicAuth(c.Username, c.Password)
	return nil
}

func (c *Client) Do(req *retryablehttp.Request, v interface{}) (*contracts.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	switch os := resp.StatusCode; os {
	case http.StatusInternalServerError:
		return nil, errors.New("server error")
	case http.StatusUnauthorized:
		return nil, errors.New("unauthorized to access this resource")
	}

	defer func() {
		err = resp.Body.Close()
	}()

	defer func() {
		_, err = io.Copy(io.Discard, resp.Body)
	}()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return &contracts.Response{Response: resp}, nil
}
