package devices

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rubemlrm/go-shelly/shelly/gen1/contracts/mocks"
	transport "github.com/rubemlrm/go-shelly/shelly/gen1/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func SetupRestClient(t *testing.T) (*http.ServeMux, *transport.Client) {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	// client is the client being tested.

	client, err := transport.NewRestClient(transport.ClientOptions{Hostname: server.URL})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return mux, client
}

func TestGetShelly(t *testing.T) {
	type test struct {
		title           string
		want            *BaseShellyResponse
		error           error
		wantError       bool
		wantClientError bool
		fixture         string
	}

	tests := []test{
		{
			title: "Testing Shelly response without error",
			want: &BaseShellyResponse{
				Type:         "SHSW-21",
				Mac:          "5ECF7F1632E8",
				Auth:         true,
				Fw:           "20161223-111304/master@2bc16496",
				LongId:       1,
				Discoverable: true,
			},
			error:           nil,
			wantError:       false,
			wantClientError: false,
			fixture:         "get_shelly.json",
		},
		{
			title:           "Testing Shelly response with error",
			want:            &BaseShellyResponse{},
			error:           nil,
			wantError:       true,
			wantClientError: false,
			fixture:         "get_shelly_error.json",
		},
		{
			title:           "Testing Shelly response with error",
			want:            &BaseShellyResponse{},
			error:           nil,
			wantError:       true,
			wantClientError: false,
			fixture:         "get_shelly_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := SetupRestClient(t)
			mux.HandleFunc("/shelly", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})

			cl := NewShellyService(client)
			resp, _, err := cl.GetShelly()
			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, resp)
			}
		})
	}
}

func TestNewRequestFailure(t *testing.T) {
	client := mocks.NewShellyClient(t)
	client.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("testing"))
	cl := NewShellyService(client)
	_, _, err := cl.GetShelly()
	assert.Error(t, err)

}

func TestGetSettings(t *testing.T) {
	type test struct {
		title     string
		want      *BaseSettingsResponse
		client    *transport.Client
		error     error
		wantError bool
		fixture   string
	}

	tests := []test{
		{
			title: "Testing Shelly response without error",
			want: &BaseSettingsResponse{
				PinCode: "123456",
			},
			client:    &transport.Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_settings.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseSettingsResponse{},
			client:    &transport.Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_settings_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := SetupRestClient(t)
			cl := NewShellyService(client)
			mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, _, err := cl.GetSettings()
			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.PinCode, resp.PinCode)
			}
		})
	}
}

func TestGetSettingsNewRequestFailure(t *testing.T) {
	client := mocks.NewShellyClient(t)
	client.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("testing"))
	cl := NewShellyService(client)
	_, _, err := cl.GetSettings()
	assert.Error(t, err)

}

func TestGetOta(t *testing.T) {
	type test struct {
		title     string
		want      *BaseOtaResponse
		client    *transport.Client
		error     error
		wantError bool
		fixture   string
	}

	tests := []test{
		{
			title: "Testing Shelly response without error",
			want: &BaseOtaResponse{
				Status:     "idle",
				HasUpdate:  false,
				NewVersion: "2",
				OldVersion: "1",
			},
			client:    &transport.Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_ota.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseOtaResponse{},
			client:    &transport.Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_ota_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := SetupRestClient(t)
			cl := NewShellyService(client)
			mux.HandleFunc("/ota", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, _, err := cl.GetOta()
			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, resp)
			}
		})
	}
}

func TestGetOtaNewRequestFailure(t *testing.T) {
	client := mocks.NewShellyClient(t)
	client.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("testing"))
	cl := NewShellyService(client)
	_, _, err := cl.GetOta()
	assert.Error(t, err)
}

func TestGetOtaCheck(t *testing.T) {
	type test struct {
		title     string
		want      *BaseOtaCheck
		client    *transport.Client
		error     error
		wantError bool
		fixture   string
	}

	tests := []test{
		{
			title: "Testing Shelly response without error",
			want: &BaseOtaCheck{
				Status: "ok",
			},
			client:    &transport.Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_ota_check.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseOtaCheck{},
			client:    &transport.Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_ota_check_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := SetupRestClient(t)
			cl := NewShellyService(client)
			mux.HandleFunc("/ota/check", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, _, err := cl.GetOtaCheck()
			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, resp)
			}
		})
	}
}

func TestGetOtaCheckNewRequestFailure(t *testing.T) {
	client := mocks.NewShellyClient(t)
	client.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("testing"))
	cl := NewShellyService(client)
	_, _, err := cl.GetOtaCheck()
	assert.Error(t, err)
}

func TestGetWifiScan(t *testing.T) {
	type test struct {
		title     string
		want      *BaseWifiScan
		client    *transport.Client
		error     error
		wantError bool
		fixture   string
	}

	tests := []test{
		{
			title: "Testing Shelly response without error",
			want: &BaseWifiScan{
				Wifiscan: "done",
				Results:  []BaseWifiScanResults{},
			},
			client:    &transport.Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_wifiscan.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseWifiScan{},
			client:    &transport.Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_wifiscan_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := SetupRestClient(t)
			cl := NewShellyService(client)
			mux.HandleFunc("/wifiscan", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, _, err := cl.GetWifiScan()
			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, resp)
			}
		})
	}
}

func TesGetWifiScanCheckNewRequestFailure(t *testing.T) {
	client := mocks.NewShellyClient(t)
	client.On("NewRequest", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("testing"))
	cl := NewShellyService(client)
	_, _, err := cl.GetWifiScan()
	assert.Error(t, err)
}

func fixture(path string) string {
	b, err := os.ReadFile("./testdata/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
