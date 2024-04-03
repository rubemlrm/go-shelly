package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*http.ServeMux, *Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	// client is the Gitlab client being tested.
	client, err := NewClient("http://localhost")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, client
}

func TestNewClient(t *testing.T) {
	type test struct {
		title     string
		wantError bool
		input     string
	}

	tests := []test{
		{
			title:     "Fail on wrong hostname",
			wantError: true,
			input:     "tes-12%wq+2",
		},
		{
			title:     "Success with correct hostname",
			wantError: false,
			input:     "localhost",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			c, err := NewClient(tc.input)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, c.BaseURL.String())
			}
		})
	}
}

func TestNewBasicAuthClient(t *testing.T) {
	type config struct {
		username    string
		password    string
		hostname    string
		requireAuth bool
	}
	type test struct {
		title     string
		wantError bool
		input     config
	}
	tests := []test{
		{
			title:     "fail on wrong hostname",
			wantError: true,
			input: config{
				username:    "",
				password:    "",
				hostname:    "tes-12%wq+2",
				requireAuth: true,
			},
		},
		{
			title:     "Valid hostname and user configuration",
			wantError: false,
			input: config{
				username:    faker.Username(),
				password:    faker.Password(),
				hostname:    faker.URL(),
				requireAuth: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			c, err := NewBasicAuthClient(tc.input.hostname, tc.input.username, tc.input.password, tc.input.requireAuth)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input.hostname, c.BaseURL.String())
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
		t.Run("", func(t *testing.T) {
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
