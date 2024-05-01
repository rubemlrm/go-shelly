package gen1

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
)

type ClientProxy interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

// Client act's as the entry object for sdk
type Client struct {
	BaseURL       *url.URL
	client        ClientProxy
	username      string
	password      string
	requireAuth   bool
	ShellyService *ShellyService
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

	c.ShellyService = &ShellyService{client: c}
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

	err = c.SetAdditionalHeaders(request, reqHeaders)
	if err != nil {
		return nil, err
	}

	if c.requireAuth {
		err = c.SetBasicAuth(request)
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

func (c *Client) SetAdditionalHeaders(request *retryablehttp.Request, headers http.Header) error {
	for k, v := range headers {
		request.Header[k] = v
	}
	return nil
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

	parsedResponse := &Response{Response: resp}
	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return parsedResponse, nil
}
