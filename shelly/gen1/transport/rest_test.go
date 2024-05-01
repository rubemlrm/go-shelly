package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/rubemlrm/go-shelly/shelly/gen1/transport/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/go-faker/faker/v4"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

type MockContext struct {
	mock.Mock
	err error
}

func (m *MockContext) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
}

func (m *MockContext) Done() <-chan struct{} {
	panic("implement me")
}

func (m *MockContext) Value(key any) any {
	panic("implement me")
}

func (m *MockContext) Err() error {
	return m.err
}

func TestNewClient(t *testing.T) {
	type test struct {
		title     string
		wantError bool
		input     ClientOptions
	}

	tests := []test{
		{
			title:     "Fail on wrong hostname",
			wantError: true,
			input: ClientOptions{
				Hostname: "tes-12%wq+2",
			},
		},
		{
			title:     "Success with correct hostname",
			wantError: false,
			input: ClientOptions{
				Hostname: "localhost",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			c, err := NewRestClient(tc.input)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input.Hostname, c.baseURL.String())
			}
		})
	}
}

func TestNewBasicAuthClient(t *testing.T) {
	type test struct {
		title     string
		wantError bool
		input     ClientOptions
	}
	tests := []test{
		{
			title:     "fail on wrong hostname",
			wantError: true,
			input: ClientOptions{
				Username: "",
				Password: "",
				Hostname: "tes-12%wq+2",
			},
		},
		{
			title:     "Valid hostname and user configuration",
			wantError: false,
			input: ClientOptions{
				Username: faker.Username(),
				Password: faker.Password(),
				Hostname: faker.URL(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			c, err := NewRestBasicAuthClient(tc.input)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input.Hostname, c.baseURL.String())
			}
		})
	}
}

func TestRetryHTTPCheck(t *testing.T) {
	type test struct {
		title            string
		input            *Client
		wantError        bool
		errorMessage     error
		httpResponse     *http.Response
		contextError     error
		wantContextError bool
		wantResponse     bool
	}
	tests := []test{
		{
			title: "Context Error found",
			input: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:        true,
			errorMessage:     nil,
			httpResponse:     &http.Response{},
			wantContextError: true,
			contextError:     errors.New("testing"),
			wantResponse:     false,
		},
		{
			title: "Error passed to method",
			input: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:        true,
			errorMessage:     errors.New("username can't be empty"),
			httpResponse:     &http.Response{},
			wantContextError: false,
			contextError:     nil,
			wantResponse:     false,
		},
		{
			title: "HTTP response retrieves a 5xx error",
			input: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:    false,
			errorMessage: nil,
			httpResponse: &http.Response{
				StatusCode: 500,
			},
			wantContextError: false,
			contextError:     nil,
			wantResponse:     true,
		},
		{
			title: "HTTP response retrieves a 200 code",
			input: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:    false,
			errorMessage: nil,
			httpResponse: &http.Response{
				StatusCode: 200,
			},
			wantContextError: false,
			contextError:     nil,
			wantResponse:     false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mockCTX := &MockContext{}
			if tc.wantContextError {
				mockCTX.err = tc.contextError
			}
			flag, err := tc.input.RetryHTTPCheck(mockCTX, tc.httpResponse, tc.errorMessage)
			if tc.wantError {
				assert.NotNil(t, err)
				if tc.wantContextError {
					assert.Equal(t, tc.contextError, err)
				} else {
					assert.Equal(t, tc.errorMessage, err)
				}
			}
			assert.Equal(t, tc.wantResponse, flag)
		})
	}
}

func TestNewRequest(t *testing.T) {

	type Opts struct {
		Param string `url:"q"`
	}

	type test struct {
		title     string
		method    string
		endpoint  string
		client    *Client
		wantError bool
		hasAuth   bool
		error     error
		opts      *Opts
	}
	url, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	tests := []test{
		{
			title:    "Request created with success",
			method:   http.MethodPost,
			endpoint: "/random",
			client: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
				baseURL:  url,
			},
			wantError: false,
			error:     nil,
			hasAuth:   false,
		},
		{
			title:    "Request with auth created with success",
			method:   http.MethodPost,
			endpoint: "/random",
			client: &Client{
				username:     faker.Username(),
				password:     faker.Password(),
				client:       &retryablehttp.Client{},
				baseURL:      url,
				requiresAuth: true,
			},
			wantError: false,
			error:     nil,
			hasAuth:   true,
		},
		{
			title:    "Request created with url parameters and with success",
			method:   http.MethodGet,
			endpoint: "/random",
			client: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
				baseURL:  url,
			},
			wantError: false,
			error:     nil,
			hasAuth:   false,
			opts:      &Opts{"testing"},
		},
		{
			title:    "Fail because of url parsing",
			method:   http.MethodGet,
			endpoint: "/random",
			client: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
				baseURL:  url,
			},
			wantError: false,
			error:     nil,
			hasAuth:   false,
			opts:      &Opts{"testing"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			req, err := tc.client.NewRequest(tc.method, tc.endpoint, tc.opts)
			assert.NoError(t, err)
			assert.NotNil(t, req)
			assert.Equal(t, tc.method, req.Request.Method)
			assert.Equal(t, []string{"application/json"}, req.Request.Header["Accept"])
			if tc.opts != nil {
				assert.Contains(t, req.Request.URL.RawQuery, "testing")
			}
			if tc.hasAuth {
				assert.NotNil(t, req.Request.Header["Authorization"])
			}
		})
	}
}

func TestSetBasicAuth(t *testing.T) {
	type test struct {
		title        string
		input        *Client
		wantError    bool
		errorMessage string
	}

	tests := []test{
		{
			title: "Username can't be empty",
			input: &Client{
				username: "",
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:    true,
			errorMessage: "username can't be empty",
		},
		{
			title: "Password can't be empty",
			input: &Client{
				username: faker.Username(),
				password: "",
				client:   &retryablehttp.Client{},
			},
			wantError:    true,
			errorMessage: "password can't be empty",
		},
		{
			title: "Auth set with success",
			input: &Client{
				username: faker.Username(),
				password: faker.Password(),
				client:   &retryablehttp.Client{},
			},
			wantError:    false,
			errorMessage: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			req, err := retryablehttp.NewRequest("GET", faker.URL(), nil)
			assert.NoError(t, err)
			err = tc.input.SetBasicAuth(req)
			if tc.wantError {
				assert.Errorf(t, err, tc.errorMessage)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDo(t *testing.T) {
	type mockClientReturn struct {
		response *http.Response
		error    error
	}
	type test struct {
		title            string
		client           *Client
		wantError        bool
		errorMessage     error
		httpResponse     *http.Response
		wantResponse     *http.Response
		mockClientReturn mockClientReturn
	}
	type resp struct {
		Title string `json:"title"`
	}

	dummyData := resp{
		Title: "testing",
	}

	responseBody, err := json.Marshal(dummyData)
	assert.NoError(t, err)

	tests := []test{
		{
			title: "Test response error",
			client: &Client{
				username: "",
				password: faker.Password(),
			},
			mockClientReturn: mockClientReturn{
				response: nil,
				error:    fmt.Errorf("testing"),
			},
			wantError: true,
		},
		{
			title: "Test response status error code",
			client: &Client{
				username: "",
				password: faker.Password(),
			},
			mockClientReturn: mockClientReturn{
				response: &http.Response{StatusCode: http.StatusInternalServerError},
				error:    nil,
			},
			wantError: true,
		},
		{
			title: "Test response status error code",
			client: &Client{
				username: faker.Username(),
				password: faker.Password(),
			},
			mockClientReturn: mockClientReturn{
				response: &http.Response{StatusCode: http.StatusUnauthorized},
				error:    nil,
			},
			wantError: true,
		},
		{
			title: "Test Failed response decode",
			client: &Client{
				username: "",
				password: faker.Password(),
			},
			mockClientReturn: mockClientReturn{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("hello")),
				},
				error: nil,
			},
			wantError: true,
		},
		{
			title: "Test Response output",
			client: &Client{
				username: "",
				password: faker.Password(),
			},
			mockClientReturn: mockClientReturn{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(io.Reader(bytes.NewReader(responseBody))),
				},
				error: nil,
			},
			wantError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mockClient := mocks.NewClientProxy(t)
			mockClient.On("Do", mock.Anything).Return(tc.mockClientReturn.response, tc.mockClientReturn.error)
			tc.client.client = mockClient
			req := &retryablehttp.Request{}
			v := &resp{}
			response, err := tc.client.Do(req, v)
			if tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, response.StatusCode, http.StatusOK)
			}
		})
	}
}
