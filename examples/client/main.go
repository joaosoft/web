package main

import (
	"fmt"
	"web"

	"github.com/joaosoft/auth-types/jwt"
)

func main() {
	// create a new client
	c, err := web.NewClient(web.WithClientMultiAttachmentMode(web.MultiAttachmentModeBoundary))
	if err != nil {
		panic(err)
	}

	requestGet(c)

	requestGetBoundary(c)

	requestAuthBasic(c)
	requestAuthJwt(c)
}

func requestGet(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/hello/joao?a=1&b=2&c=1,2,3")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
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

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", response)
}

func requestAuthBasic(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-basic")
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthBasic("joao", "ribeiro")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}

func requestAuthJwt(c *web.Client) {
	request, err := c.NewRequest(web.MethodGet, "localhost:9001/auth-jwt")
	if err != nil {
		panic(err)
	}

	_, err = request.WithAuthJwt(jwt.Claims{"name": "joao", "age": 30}, "bananas")
	if err != nil {
		panic(err)
	}

	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\n%d: %s\n\n", response.Status, string(response.Body))
}