package main

import shelly "github.com/rubemlrm/go-shelly/shelly"

func main() {
	client, err := shelly.NewClient("http://localhost")
	if err != nil {
		panic(err)
	}
	var rp interface{}
	req, err := client.NewRequest("GET", "/settings", rp)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req, rp)
	if err != nil {
		panic(err)
	}
	print(resp)
}
