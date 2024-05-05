package gen1

import (
	devices "github.com/rubemlrm/go-shelly/shelly/gen1/devices"
	transport "github.com/rubemlrm/go-shelly/shelly/gen1/transport"
)

type RestClient struct {
	ShellyService *devices.ShellyService
}

func NewRestClient(options transport.ClientOptions) (*RestClient, error) {
	cl, err := transport.NewRestClient(options)
	if err != nil {
		return nil, err
	}
	return &RestClient{
		ShellyService: devices.NewShellyService(cl),
	}, nil
}

func NewRestClientWithAuth(options transport.ClientOptions) (*RestClient, error) {
	cl, err := transport.NewRestBasicAuthClient(options)
	if err != nil {
		return nil, err
	}
	return &RestClient{
		ShellyService: devices.NewShellyService(cl),
	}, nil
}
