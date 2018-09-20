package main

import (
	"web/client"
)

func main() {
	// create a new server
	c, err := client.NewWebClient()
	if err != nil {
		panic(err)
	}

	request, _ := c.NewRequest()
	c.Call(request)

}
