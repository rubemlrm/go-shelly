package gen1

import (
	"net/url"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/rubemlrm/go-shelly/shelly/gen1/devices"
	transport "github.com/rubemlrm/go-shelly/shelly/gen1/transport"
	"github.com/stretchr/testify/assert"
)

func TestNewRestClient(t *testing.T) {
	type args struct {
		options transport.ClientOptions
	}
	url, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wants   RestClient
	}{
		{
			name: "Test Client creation with sucess",
			args: args{
				options: transport.ClientOptions{
					Hostname: "http://localhost",
				},
			},
			wantErr: false,
			wants: RestClient{
				ShellyService: &devices.ShellyService{
					Client: &transport.Client{
						BaseURL: url,
					},
				},
			},
		},
		{
			name: "Test Client creation with error",
			args: args{
				options: transport.ClientOptions{
					Hostname: "http://»%@ 2 2.com",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRestClient(tt.args.options)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				v := client.ShellyService.Client.(*transport.Client)
				assert.Equal(t, v.BaseURL, tt.wants.ShellyService.Client.(*transport.Client).BaseURL)

			}
		})
	}
}

func TestNewRestClientWithAuth(t *testing.T) {
	type args struct {
		options transport.ClientOptions
	}
	username := faker.Username()
	password := faker.Password()
	url, err := url.Parse("http://localhost")
	assert.NoError(t, err)

	tests := []struct {
		name    string
		args    args
		wantErr bool
		wants   RestClient
	}{
		{
			name: "Test Client creation with sucess",
			args: args{
				options: transport.ClientOptions{
					Hostname: "http://localhost",
					Username: username,
					Password: password,
				},
			},
			wantErr: false,
			wants: RestClient{
				ShellyService: &devices.ShellyService{
					Client: &transport.Client{
						BaseURL:  url,
						Username: username,
						Password: password,
					},
				},
			},
		},
		{
			name: "Test Client creation with error",
			args: args{
				options: transport.ClientOptions{
					Hostname: "http://»%@ 2 2.com",
					Username: faker.Username(),
				},
			},
			wantErr: true,
			wants: RestClient{
				ShellyService: &devices.ShellyService{
					Client: &transport.Client{
						BaseURL: url,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRestClientWithAuth(tt.args.options)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)

			}
		})
	}
}
