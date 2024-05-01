package gen1

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := setup(t)
			mux.HandleFunc("/shelly", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, err := client.ShellyService.GetShelly()
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

func TestGetSettings(t *testing.T) {
	type test struct {
		title     string
		want      *BaseSettingsResponse
		client    *Client
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
			client:    &Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_settings.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseSettingsResponse{},
			client:    &Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_settings_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := setup(t)
			mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, err := client.ShellyService.GetSettings()
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

func TestGetOta(t *testing.T) {
	type test struct {
		title     string
		want      *BaseOtaResponse
		client    *Client
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
			client:    &Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_ota.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseOtaResponse{},
			client:    &Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_ota_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := setup(t)
			mux.HandleFunc("/ota", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, err := client.ShellyService.GetOta()
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

func TestGetOtaCheck(t *testing.T) {
	type test struct {
		title     string
		want      *BaseOtaCheck
		client    *Client
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
			client:    &Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_ota_check.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseOtaCheck{},
			client:    &Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_ota_check_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := setup(t)
			mux.HandleFunc("/ota/check", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, err := client.ShellyService.GetOtaCheck()
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

func TestGetWifiScan(t *testing.T) {
	type test struct {
		title     string
		want      *BaseWifiScan
		client    *Client
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
			client:    &Client{},
			error:     nil,
			wantError: false,
			fixture:   "get_wifiscan.json",
		},
		{
			title:     "Testing Shelly response with error",
			want:      &BaseWifiScan{},
			client:    &Client{},
			error:     nil,
			wantError: true,
			fixture:   "get_wifiscan_error.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			mux, client := setup(t)
			mux.HandleFunc("/wifiscan", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				fmt.Fprint(w, fixture(tc.fixture))
			})
			resp, err := client.ShellyService.GetWifiScan()
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

func fixture(path string) string {
	b, err := os.ReadFile("./testdata/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
