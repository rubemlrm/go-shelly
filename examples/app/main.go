package main

import (
	"fmt"
	"os"

	gen1 "github.com/rubemlrm/go-shelly/shelly/gen1"
	transport "github.com/rubemlrm/go-shelly/shelly/gen1/transport"
)

func main() {
	opts := transport.ClientOptions{
		Hostname: os.Getenv("HOST"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
	client, err := gen1.NewRestClientWithAuth(opts)
	if err != nil {
		panic(err)
	}
	req, err := client.ShellyService.GetSettings()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", req)
	req2, err := client.ShellyService.GetOta()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", req2)

	req3, err := client.ShellyService.GetOtaCheck()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", req3)

	req4, err := client.ShellyService.GetWifiScan()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", req4)
}
