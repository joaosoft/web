package main

import (
	"fmt"
	"web"
)

func main() {
	// create a new client
	c, err := web.NewClient(web.WithClientMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)

	requestGetBoundary(c)

}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestGetBoundary(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao/download")
	if err != nil {
		panic(err)
	}

	response, err := c.Send(request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}