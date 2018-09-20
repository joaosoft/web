package main

import (
	"fmt"
	"web/client"
	"web/common"
)

func main() {
	// create a new client
	c, err := client.NewClient(client.WithMultiAttachmentMode(common.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)
	requestGet(c)
	requestGet(c)

	requestGetBoundary(c)

}

func requestGet(c *client.Client) {
	request, err := c.NewRequest(common.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestGetBoundary(c *client.Client) {
	request, err := c.NewRequest(common.MethodGet, "localhost:9001/hello/joao/download")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}
