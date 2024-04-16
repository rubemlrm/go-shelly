package main

import (
	"fmt"
	"os"

	shelly "github.com/rubemlrm/go-shelly/shelly/gen1"
)

func main() {
	client, err := shelly.NewBasicAuthClient(os.Getenv("HOST"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), true)
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
